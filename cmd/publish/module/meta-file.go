package module

import (
	"fmt"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/module"
)

type MetaFile struct {
	XTerraformGet string `json:"X-Terraform-Get"`
}

func createMetaFiles(module module.Module, file module.RepositoryFile) error {
	for _, d := range file.Versions {
		version := d.Version

		filePath := getMetaFilePath(module, version)
		fileContent := MetaFile{XTerraformGet: getXTerraformGet(module, version)}

		err := files.WriteToFile(filePath, fileContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func getMetaFilePath(module module.Module, version string) string {
	return fmt.Sprintf("dist/v1/modules/%s/%s/%s/%s/download", module.Namespace, module.Name, module.System, internal.TrimTagPrefix(version))
}

func getXTerraformGet(mod module.Module, version string) string {
	return fmt.Sprintf("git::%s?ref=%s", module.GetRepositoryUrl(mod), version)
}
