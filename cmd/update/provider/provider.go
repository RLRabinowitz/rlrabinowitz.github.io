package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"golang.org/x/mod/semver"
	"log"
	"net/http"
	providerInitialize "rlrabinowitz.github.io/cmd/initialize/provider"
	"rlrabinowitz.github.io/internal"

	"os"
	"path"
	"rlrabinowitz.github.io/internal/github"
	"rlrabinowitz.github.io/internal/provider"
	"strings"
)

func Update(pathToFile string, experimental bool) error {
	// TODO Better checks for this?
	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	p := provider.Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	if !experimental {
		shouldUpdate, err := shouldUpdateByTags(p, pathToFile)
		if err != nil {
			return err
		}

		if shouldUpdate {
			return providerInitialize.Initialize(pathToFile)
		}
	} else {
		tagsShouldUpdate, err := shouldUpdateByTags(p, pathToFile)
		if err != nil {
			log.Printf("Experimental: Could not check releases via tags for %s - %s", pathToFile, err)
			return err
		}

		rssShouldUpdate, err := shouldUpdateByRss(p, pathToFile)
		if err != nil {
			log.Printf("Experimental: Could not check releases via RSS for %s - %s", pathToFile, err)
			return err
		}

		if tagsShouldUpdate && rssShouldUpdate {
			log.Printf("Experimental Result: %s was found by both methods", pathToFile)
		} else if tagsShouldUpdate {
			log.Printf("Experimental Result: %s was found only by tags", pathToFile)
		} else if rssShouldUpdate {
			log.Printf("Experimental Result: %s was found only by RSS", pathToFile)
		} else {
			log.Printf("Experimental: %s does not require update", pathToFile)
		}
	}

	//highestSemverTag, err := getHighestSemverTag(p)
	//if err != nil {
	//	return err
	//}
	//
	//fileContent, err := getProviderFileContent(pathToFile)
	//if err != nil {
	//	return err
	//}
	//
	//for _, v := range fileContent.Versions {
	//	versionWithPrefix := fmt.Sprintf("v%s", v.Version)
	//	if versionWithPrefix == highestSemverTag {
	//		log.Printf("Found latest tag %s in the repository file %s, nothing to update...", highestSemverTag, pathToFile)
	//		return nil
	//	}
	//}
	//
	//log.Printf("Could not find latest tag %s in the repository file %s, updating the file...", highestSemverTag, pathToFile)
	//
	//return providerInitialize.Initialize(pathToFile)
	return nil
}

func shouldUpdateByTags(p provider.Provider, pathToFile string) (bool, error) {
	highestSemverTag, err := getHighestSemverTag(p)
	if err != nil {
		return false, err
	}

	fileContent, err := getProviderFileContent(pathToFile)
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == highestSemverTag {
			log.Printf("Found latest tag %s in the repository file %s, nothing to update...", highestSemverTag, pathToFile)
			return false, nil
		}
	}

	log.Printf("Could not find latest tag %s in the repository file %s, updating the file...", highestSemverTag, pathToFile)
	return true, nil
}

func getRssFeedAlternative(url string) (*gofeed.Feed, error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gofeed/1.0")

	token := os.Getenv("GH_TOKEN")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("got error %d", resp.StatusCode)
	}

	return gofeed.NewParser().Parse(resp.Body)
}

func shouldUpdateByRss(p provider.Provider, pathToFile string) (bool, error) {
	//fp := gofeed.NewParser()
	rssUrl := getRssUrl(p)
	//feed, err := fp.ParseURL(rssUrl)
	feed, err := getRssFeedAlternative(rssUrl)
	if err != nil {
		return false, err
	}

	if len(feed.Items) < 1 {
		log.Printf("No releases found in RSS feed %s", rssUrl)
	}

	log.Printf("Found %d releases in RSS feed %s", len(feed.Items), rssUrl)
	latestRelease := feed.Items[0]
	tagName := latestRelease.Title

	fileContent, err := getProviderFileContent(pathToFile)
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == fmt.Sprintf("v%s", internal.TrimTagPrefix(tagName)) {
			log.Printf("Found latest tag %s in the repository file %s, nothing to update...", tagName, pathToFile)
			return false, nil
		}
	}

	log.Printf("Could not find latest tag %s in the repository file %s, updating the file...", tagName, pathToFile)
	return true, nil
}

func getRssUrl(p provider.Provider) string {
	repositoryUrl := provider.GetRepositoryUrl(p)
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}

// TODO Commonize
func getProviderFileContent(path string) (provider.RepositoryFile, error) {
	res, _ := os.ReadFile(path)

	var fileData provider.RepositoryFile

	err := json.Unmarshal(res, &fileData)
	// TODO better error handling
	if err != nil {
		return provider.RepositoryFile{}, err
	}

	return fileData, nil
}

func getSemverTags(p provider.Provider) ([]string, error) {
	repositoryUrl := provider.GetRepositoryUrl(p)
	tags, err := github.GetTags(repositoryUrl)
	if err != nil {
		return nil, err
	}
	var semverTags []string

	for _, tag := range tags {
		if semver.IsValid(tag) {
			semverTags = append(semverTags, tag)
		}
	}

	semver.Sort(semverTags)

	return semverTags, nil
}

func getHighestSemverTag(p provider.Provider) (string, error) {
	semverTags, err := getSemverTags(p)
	if err != nil {
		return "", err
	}

	if len(semverTags) < 1 {
		return "", fmt.Errorf("no semver tags found in repository for provider %s/%s", p.Namespace, p.ProviderName)
	}

	return semverTags[len(semverTags)-1], nil
}
