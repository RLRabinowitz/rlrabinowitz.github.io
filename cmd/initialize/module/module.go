package module

import (
	"path"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/github"
	"rlrabinowitz.github.io/internal/module"
	"strings"
)

func Initialize(pathToFile string) error {
	fileName := path.Base(pathToFile)
	system := strings.TrimSuffix(fileName, path.Ext(fileName))
	name := path.Base(path.Dir(pathToFile))
	namespace := path.Base(path.Dir(path.Dir(pathToFile)))
	mod := module.Module{
		Namespace: namespace,
		Name:      name,
		System:    system,
	}

	fileContent, err := toRepositoryFileData(mod)
	if err != nil {
		return err
	}

	return files.WriteToFile(pathToFile, fileContent)
}

func getModuleTags(mod module.Module) ([]string, error) {
	return github.GetTags(module.GetRepositoryUrl(mod))
}

func toRepositoryFileData(mod module.Module) (*module.RepositoryFile, error) {
	tags, err := getModuleTags(mod)
	if err != nil {
		return nil, err
	}

	var versions = make([]module.Version, 0)
	for _, t := range tags {
		versions = append(versions, module.Version{Version: t})
	}

	return &module.RepositoryFile{Versions: versions}, nil
}
