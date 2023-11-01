package provider

import (
	"encoding/json"
	"fmt"
	"golang.org/x/mod/semver"
	"log"
	"os"
	"path"
	providerInitialize "rlrabinowitz.github.io/cmd/initialize/provider"
	"rlrabinowitz.github.io/internal/github"
	"rlrabinowitz.github.io/internal/provider"
	"strings"
)

func Update(pathToFile string) error {
	// TODO Better checks for this?
	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	p := provider.Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	highestSemverTag, err := getHighestSemverTag(p)
	if err != nil {
		return err
	}

	fileContent, err := getProviderFileContent(pathToFile)
	if err != nil {
		return err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == highestSemverTag {
			log.Printf("Found latest tag %s in the repository file %s, nothing to update...", highestSemverTag, pathToFile)
			return nil
		}
	}

	log.Printf("Could not find latest tag %s in the repository file %s, updating the file...", highestSemverTag, pathToFile)

	return providerInitialize.Initialize(pathToFile)
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
		return "", nil
	}

	if len(semverTags) < 1 {
		return "", fmt.Errorf("no semver tags found in repository for provider %s/%s", p.Namespace, p.ProviderName)
	}

	return semverTags[len(semverTags)-1], nil
}
