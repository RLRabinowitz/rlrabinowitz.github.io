package module

import (
	"encoding/json"
	"os"
	"path"
	"rlrabinowitz.github.io/internal/module"
	"strings"
)

func Update(pathToFile string) error {
	fileName := path.Base(pathToFile)
	system := strings.TrimSuffix(fileName, path.Ext(fileName))
	name := path.Base(path.Dir(pathToFile))
	namespace := path.Base(path.Dir(path.Dir(pathToFile)))
	mod := module.Module{
		Namespace: namespace,
		Name:      name,
		System:    system,
	}

	fileContent, err := getModuleFileContent(pathToFile)
	if err != nil {
		return err
	}

}

// TODO Commonize
func getModuleFileContent(pathToFile string) (module.RepositoryFile, error) {
	res, _ := os.ReadFile(pathToFile)

	var fileData module.RepositoryFile

	err := json.Unmarshal(res, &fileData)
	// TODO better error handling
	if err != nil {
		return module.RepositoryFile{}, err
	}

	return fileData, nil
}
