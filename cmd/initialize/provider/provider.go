package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shurcooL/githubv4"
	"os"
	"path"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/github"
	"rlrabinowitz.github.io/internal/provider"
	"strings"
)

// TODO better name or better arg
func Initialize(pathToFile string) error {
	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	p := provider.Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	// TODO move context back?
	ctx := context.Background()
	client := github.NewGitHubClient(ctx, os.Getenv("GH_TOKEN"))
	fileContent, err := getInputData(ctx, client, p)
	if err != nil {
		return err
	}

	marshalledJson, err := json.Marshal(fileContent)
	if err != nil {
		return err
	}

	return files.WriteToFile(pathToFile, marshalledJson)
}

func getInputData(ctx context.Context, ghClient *githubv4.Client, p provider.Provider) (*provider.RepositoryFile, error) {
	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.Namespace, p.ProviderName)
	if err != nil {
		return nil, err
	}

	versions := make([]provider.Version, 0)
	for _, r := range releases {
		var shaSumsArtifact *github.ReleaseAsset
		var shaSumsSignatureArtifact *github.ReleaseAsset

		var targets = make([]provider.Target, 0)
		for _, asset := range r.ReleaseAssets.Nodes {
			if platform := github.ExtractPlatformFromFilename(asset.Name); platform != nil {
				shaSum, err := github.GetShaSum(ctx, asset.DownloadURL, asset.Name)
				if err != nil {
					return nil, err
				}
				targets = append(targets, provider.Target{
					OS:          platform.OS,
					Arch:        platform.Arch,
					Filename:    asset.Name,
					DownloadURL: asset.DownloadURL,
					SHASum:      shaSum,
				})
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS") {
				shaSumsArtifact = &asset
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS.sig") {
				shaSumsSignatureArtifact = &asset
			}
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

	return &provider.RepositoryFile{Versions: versions}, nil
}

func getRepositoryUrl(provider provider.Provider) string {
	return fmt.Sprintf("https://github.com/%s/%s", provider.Namespace, provider.ProviderName)
}
