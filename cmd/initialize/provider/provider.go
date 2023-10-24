package provider

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"log"
	"os"
	"path"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/github"
	"rlrabinowitz.github.io/internal/provider"
	"strings"
	"sync"
)

type ShaSumResult struct {
	Index     int
	ShaSumMap map[string]string
	Err       error
}

// TODO better name or better arg
func Initialize(pathToFile string) error {
	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	p := provider.Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	ctx := context.Background()
	client := github.NewGitHubClient(ctx, os.Getenv("GH_TOKEN"))
	fileContent, err := toRepositoryFileData(ctx, client, p)
	if err != nil {
		return err
	}

	return files.WriteToFile(pathToFile, fileContent)
}

func toRepositoryFileData(ctx context.Context, ghClient *githubv4.Client, p provider.Provider) (*provider.RepositoryFile, error) {
	repoName := provider.GetRepositoryName(p)
	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.Namespace, repoName)
	if err != nil {
		return nil, err
	}

	versions := make([]provider.Version, 0)
	for _, r := range releases {
		var shaSumsArtifact github.ReleaseAsset
		var shaSumsSignatureArtifact github.ReleaseAsset

		var targets = make([]provider.Target, 0)
		for _, asset := range r.ReleaseAssets.Nodes {
			if platform := github.ExtractPlatformFromFilename(asset.Name); platform != nil {
				if err != nil {
					return nil, err
				}
				targets = append(targets, provider.Target{
					OS:          platform.OS,
					Arch:        platform.Arch,
					Filename:    asset.Name,
					DownloadURL: asset.DownloadURL,
				})
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS") {
				shaSumsArtifact = asset
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS.sig") {
				shaSumsSignatureArtifact = asset
			}
		}
		if len(targets) == 0 {
			log.Printf("could not find artifacts in release of provider %s version %s, skipping...", p.ProviderName, r.TagName)
			continue
		}
		if (shaSumsArtifact == github.ReleaseAsset{}) {
			return nil, fmt.Errorf("could not SHASUMS artifact for provider %s version %s", p.ProviderName, r.TagName)
		}
		if (shaSumsSignatureArtifact == github.ReleaseAsset{}) {
			return nil, fmt.Errorf("could not SHASUMS signature artifact for provider %s version %s", p.ProviderName, r.TagName)
		}

		versions = append(versions, provider.Version{
			Version:             internal.TrimTagPrefix(r.TagName),
			Protocols:           []string{"5.0"}, // TODO fetch manifest
			SHASumsURL:          shaSumsArtifact.DownloadURL,
			SHASumsSignatureURL: shaSumsSignatureArtifact.DownloadURL,
			Repository:          getRepositoryUrl(p), // TODO move the repository setting
			// TODO actually - how to know the repository initially???
			Targets: targets,
		})
	}

	versions, err = enrichWithShaSums(ctx, versions)
	if err != nil {
		return nil, err
	}

	return &provider.RepositoryFile{Versions: versions}, nil
}

func enrichWithShaSums(ctx context.Context, versions []provider.Version) ([]provider.Version, error) {
	versionsCopy := versions

	shaSumCh := make(chan ShaSumResult, len(versionsCopy))

	var wg sync.WaitGroup
	for i, v := range versionsCopy {
		wg.Add(1)

		go func(v provider.Version, i int) {
			defer wg.Done()
			shaMap, err := github.GetShaSums(ctx, v.SHASumsURL)
			shaSumCh <- ShaSumResult{
				Index:     i,
				ShaSumMap: shaMap,
				Err:       err,
			}
		}(v, i)
	}

	wg.Wait()
	close(shaSumCh)

	for sr := range shaSumCh {
		if sr.Err != nil {
			return nil, fmt.Errorf("failed to find SHA of artifact: %w", sr.Err)
		}

		for i, t := range versionsCopy[sr.Index].Targets {
			shaSum := sr.ShaSumMap[t.Filename]
			versionsCopy[sr.Index].Targets[i].SHASum = shaSum
		}
	}

	return versionsCopy, nil
}

func getRepositoryUrl(provider provider.Provider) string {
	return fmt.Sprintf("https://github.com/%s/%s", provider.Namespace, provider.ProviderName)
}
