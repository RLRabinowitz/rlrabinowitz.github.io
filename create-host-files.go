package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// TODO Validate module tags (and provider tags) make sense

type ModuleInputVersion struct {
	Version string `json:"version"`
}

type ModuleInputVersionFile struct {
	Versions []ModuleInputVersion `json:"versions"`
}

type InputVersion struct {
	Version             string   `json:"version"`
	Protocols           []string `json:"protocols"`
	SHASumsURL          string   `json:"shasums_url"`
	SHASumsSignatureURL string   `json:"shasums_signature_url"`
	Repository          string   `json:"repository"`
	Targets             []struct {
		OS          string `json:"os"`           // The operating system for which the provider is built.
		Arch        string `json:"arch"`         // The architecture for which the provider is built.
		Filename    string `json:"filename"`     // The filename of the provider binary.
		DownloadURL string `json:"download_url"` // The direct URL to download the provider binary.
		SHASum      string `json:"shasum"`       // The SHA checksum of the provider binary.
	} `json:"targets"`
}

type GPGPublicKeys struct {
	KeyID      string `json:"key_id"`      // The ID of the GPG key.
	ASCIIArmor string `json:"ascii_armor"` // The ASCII armored representation of the GPG public key.
}

type InputVersionFile struct {
	Versions []InputVersion `json:"versions"`
}

type OutputPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// TODO: Handle warnings
type OutputVersion struct {
	Version   string           `json:"version"`   // The version number of the provider.
	Protocols []string         `json:"protocols"` // The protocol versions the provider supports.
	Platforms []OutputPlatform `json:"platforms"` // A list of platforms for which this provider version is available.
}

type OutputVersionsFile struct {
	Versions []OutputVersion `json:"versions"`
	Warnings []string        `json:"warnings,omitempty"`
}

type OutputModuleVersionsFile struct {
	Modules []ModuleInputVersionFile `json:"modules"`
}

// TODO Removal of version, and removal of entire provider

type SigningKeys struct {
	GPGPublicKeys []GPGPublicKeys `json:"gpg_public_keys"` // A list of GPG public keys.
}

type OutputModuleVersionMetaFile struct {
	XTerraformGet string `json:"X-Terraform-Get"`
}

type OutputVersionMetaFile struct {
	Protocols           []string    `json:"protocols"`             // The protocol versions the provider supports.
	OS                  string      `json:"os"`                    // The operating system for which the provider is built.
	Arch                string      `json:"arch"`                  // The architecture for which the provider is built.
	Filename            string      `json:"filename"`              // The filename of the provider binary.
	DownloadURL         string      `json:"download_url"`          // The direct URL to download the provider binary.
	SHASumsURL          string      `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string      `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	SHASum              string      `json:"shasum"`                // The SHA checksum of the provider binary.
	SigningKeys         SigningKeys `json:"signing_keys"`          // The signing keys used for this provider version.
}

type Provider struct {
	ProviderName string
	Namespace    string
}

type Module struct {
	Namespace string
	Name      string
	System    string
}

func main() {
	filePaths := getFilePathsToMigrate()
	for _, filePath := range filePaths {
		// TODO Range variables
		// TODO Parallelism
		if isProviderPath(filePath) {
			err := runForProviderFile(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			err := runForModuleFile(filePath)
			if err != nil {
				panic(err)
			}
		} // TODO Validate path is either provider or module, that amount of parts make sense
	}
}

func getFilePathsToMigrate() []string {
	if len(os.Args) != 2 {
		panic("The fuck, missing arguments") // TODO language (and panic)
	}
	return strings.Split(os.Args[1], ",")
}

func isProviderPath(filePath string) bool {
	pathParts := strings.Split(filePath, "/")
	return pathParts[0] == "providers" // TODO constant?
}

func runForModuleFile(pathToFile string) error {
	fileName := path.Base(pathToFile)
	system := strings.TrimSuffix(fileName, path.Ext(fileName))
	name := path.Base(path.Dir(pathToFile))
	namespace := path.Base(path.Dir(path.Dir(pathToFile)))
	module := Module{
		Namespace: namespace,
		Name:      name,
		System:    system,
	}

	fileContent, err := getModuleFileContent(pathToFile)
	if err != nil {
		return err
	}

	// TODO - Validate file? Other than JSON Unmarshal? Required Fields? Validation for fields?

	err = createModuleVersionsFile(module, fileContent)
	if err != nil {
		return err
	}

	err = createModuleDownloadFiles(module, fileContent)
	if err != nil {
		return err
	}

	return nil

}

func getPathToModuleVersionsFile(module Module) string {
	return fmt.Sprintf("dist/v1/modules/%s/%s/%s/versions", module.Namespace, module.Name, module.System)
}

func trimModuleTagPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}

func getPathToModuleDownloadFile(module Module, version string) string {
	return fmt.Sprintf("dist/v1/modules/%s/%s/%s/%s/download", module.Namespace, module.Name, module.System, trimModuleTagPrefix(version))
}

func createModuleVersionsFile(module Module, file ModuleInputVersionFile) error {
	filePath := getPathToModuleVersionsFile(module)

	data := mapToModuleOutputVersions(file)
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

// TODO How to handle 404s and other such errors?

func getXTerraformGet(module Module, version string) string {
	repoName := fmt.Sprintf("terraform-%s-%s", module.System, module.Name) // TODO hashi?
	return fmt.Sprintf("git::https://github.com/%s/%s?ref=%s", module.Namespace, repoName, version)
}

func createModuleDownloadFiles(module Module, file ModuleInputVersionFile) error {
	for _, d := range file.Versions {
		version := d.Version

		filePath := getPathToModuleDownloadFile(module, version)
		fileContent := OutputModuleVersionMetaFile{XTerraformGet: getXTerraformGet(module, version)}

		marshalledJson, err := json.Marshal(fileContent)
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

	}

	return nil
}

func getModuleFileContent(pathToFile string) (ModuleInputVersionFile, error) {
	res, _ := os.ReadFile(pathToFile)

	var fileData ModuleInputVersionFile

	err := json.Unmarshal(res, &fileData)
	// TODO better error handling
	if err != nil {
		return ModuleInputVersionFile{}, err
	}

	return fileData, nil
}

func runForProviderFile(pathToFile string) error {
	// TODO check that this is a provider, and that the path makes sense

	fileName := path.Base(pathToFile)
	providerName := strings.TrimSuffix(fileName, path.Ext(fileName))
	namespace := path.Base(path.Dir(pathToFile))
	provider := Provider{
		ProviderName: providerName,
		Namespace:    namespace,
	}

	// TODO replace Hashi

	fileContent, err := getFileContent(pathToFile)
	if err != nil {
		return err
	}

	// TODO - Validate file? Other than JSON Unmarshal? Required Fields? Validation for fields?

	err = createVersionsFile(provider, fileContent)
	if err != nil {
		return err
	}

	err = createDownloadFiles(provider, fileContent)
	if err != nil {
		return err
	}

	return nil
}

func getFileContent(path string) (InputVersionFile, error) {
	res, _ := os.ReadFile(path)

	var fileData InputVersionFile

	err := json.Unmarshal(res, &fileData)
	// TODO better error handling
	if err != nil {
		return InputVersionFile{}, err
	}

	return fileData, nil
}

func mapToModuleOutputVersions(file ModuleInputVersionFile) OutputModuleVersionsFile {
	outputVersionsFile := make([]ModuleInputVersion, len(file.Versions))
	for i, d := range file.Versions {
		outputVersionsFile[i] = ModuleInputVersion{Version: trimModuleTagPrefix(d.Version)}
	}

	return OutputModuleVersionsFile{
		Modules: []ModuleInputVersionFile{
			{
				Versions: outputVersionsFile,
			},
		},
	}
}

func mapToOutputVersions(file InputVersionFile) OutputVersionsFile {
	outputVersionsFile := make([]OutputVersion, len(file.Versions))
	for i, d := range file.Versions {
		outputPlatforms := make([]OutputPlatform, len(d.Targets))

		for ip, dp := range d.Targets {
			outputPlatforms[ip] = OutputPlatform{
				OS:   dp.OS,
				Arch: dp.Arch,
			}
		}

		outputVersionsFile[i] = OutputVersion{
			Version:   trimModuleTagPrefix(d.Version),
			Protocols: d.Protocols,
			Platforms: outputPlatforms,
		}
	}

	return OutputVersionsFile{
		Versions: outputVersionsFile,
	}
}

// TODO Do not hardcode the "dist" folder

func getPathToVersionsFile(provider Provider) string {
	return fmt.Sprintf("dist/v1/providers/%s/%s/versions", provider.Namespace, provider.ProviderName)
}

func getPathToDownloadFile(provider Provider, version string, platform OutputPlatform) string {
	return fmt.Sprintf("dist/v1/providers/%s/%s/%s/download/%s/%s", provider.Namespace, provider.ProviderName, trimModuleTagPrefix(version), platform.OS, platform.Arch)
}

func createVersionsFile(provider Provider, file InputVersionFile) error {
	filePath := getPathToVersionsFile(provider)

	data := mapToOutputVersions(file)
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

func createDownloadFiles(provider Provider, file InputVersionFile) error {
	for _, d := range file.Versions {
		version := d.Version
		for _, dp := range d.Targets {
			// TODO refactor
			platform := OutputPlatform{
				OS:   dp.OS,
				Arch: dp.Arch,
			}
			filePath := getPathToDownloadFile(provider, version, platform)

			// TODO Keys
			fileContent := OutputVersionMetaFile{
				Protocols:           d.Protocols,
				OS:                  dp.OS,
				Arch:                dp.Arch,
				Filename:            dp.Filename,
				DownloadURL:         dp.DownloadURL,
				SHASumsURL:          d.SHASumsURL,
				SHASumsSignatureURL: d.SHASumsSignatureURL,
				SHASum:              dp.SHASum,
				SigningKeys:         SigningKeys{GPGPublicKeys: make([]GPGPublicKeys, 0)},
			}

			marshalledJson, err := json.Marshal(fileContent)
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
		}
	}

	return nil
}
