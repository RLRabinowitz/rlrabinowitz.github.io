package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/provider"
)

// TODO: Handle warnings
type Version struct {
	Version   string     `json:"version"`   // The version number of the provider.
	Protocols []string   `json:"protocols"` // The protocol versions the provider supports.
	Platforms []Platform `json:"platforms"` // A list of platforms for which this provider version is available.
}

type VersionsFile struct {
	Versions []Version `json:"versions"`
	Warnings []string  `json:"warnings,omitempty"`
}

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

func createVersionsFile(provider provider.Provider, file provider.RepositoryFile) error {
	filePath := getVersionsFilePath(provider)

	data := mapToVersions(file)
	marshalledJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// TODO what if exists
	err = os.MkdirAll(path.Dir(filePath), 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, marshalledJson, 0777) // TODO Perm
	if err != nil {
		return err
	}

	return nil
}

func getVersionsFilePath(provider provider.Provider) string {
	return fmt.Sprintf("dist/v1/providers/%s/%s/versions", provider.Namespace, provider.ProviderName)
}

// TODO How to handle 404s and other such errors?

func mapToVersions(file provider.RepositoryFile) VersionsFile {
	outputVersionsFile := make([]Version, len(file.Versions))
	for i, d := range file.Versions {
		outputPlatforms := make([]Platform, len(d.Targets))

		for ip, dp := range d.Targets {
			outputPlatforms[ip] = Platform{
				OS:   dp.OS,
				Arch: dp.Arch,
			}
		}

		outputVersionsFile[i] = Version{
			Version:   internal.TrimTagPrefix(d.Version),
			Protocols: d.Protocols,
			Platforms: outputPlatforms,
		}
	}

	return VersionsFile{
		Versions: outputVersionsFile,
	}
}
