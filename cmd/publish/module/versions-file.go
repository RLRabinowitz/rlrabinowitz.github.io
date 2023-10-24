package module

import (
	"fmt"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/module"
)

type VersionsFile struct {
	Modules []module.RepositoryFile `json:"modules"`
}

func createVersionsFile(module module.Module, file module.RepositoryFile) error {
	filePath := getVersionsFilePath(module)

	data := mapToVersions(file)

	return files.WriteToFile(filePath, data)
}

func getVersionsFilePath(module module.Module) string {
	return fmt.Sprintf("dist/v1/modules/%s/%s/%s/versions", module.Namespace, module.Name, module.System)
}

func mapToVersions(file module.RepositoryFile) VersionsFile {
	outputVersionsFile := make([]module.Version, len(file.Versions))
	for i, d := range file.Versions {
		outputVersionsFile[i] = module.Version{Version: internal.TrimTagPrefix(d.Version)}
	}

	return VersionsFile{
		Modules: []module.RepositoryFile{
			{
				Versions: outputVersionsFile,
			},
		},
	}
}
