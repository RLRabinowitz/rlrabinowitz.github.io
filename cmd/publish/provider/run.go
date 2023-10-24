package provider

import (
	"encoding/json"
	"os"
	"path"
	"rlrabinowitz.github.io/internal/provider"
	"strings"
)

func Publish(pathToFile string) error {
	// TODO check that this is a provider, and that the path makes sense

	// TODO Better checks for this?
	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	provider := provider.Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	// TODO replace Hashi

	fileContent, err := getProviderFileContent(pathToFile)
	if err != nil {
		return err
	}

	// TODO - Validate file? Other than JSON Unmarshal? Required Fields? Validation for fields?

	err = createVersionsFile(provider, fileContent)
	if err != nil {
		return err
	}

	err = createMetaFiles(provider, fileContent)
	if err != nil {
		return err
	}

	return nil
}

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
