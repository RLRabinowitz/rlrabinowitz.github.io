package provider

import (
	"fmt"
	"rlrabinowitz.github.io/internal"
	"rlrabinowitz.github.io/internal/files"
	"rlrabinowitz.github.io/internal/provider"
)

// TODO Validate module tags (and provider tags) make sense

type GPGPublicKeys struct {
	KeyID      string `json:"key_id"`      // The ID of the GPG key.
	ASCIIArmor string `json:"ascii_armor"` // The ASCII armored representation of the GPG public key.
}

// TODO Removal of version, and removal of entire provider

type SigningKeys struct {
	GPGPublicKeys []GPGPublicKeys `json:"gpg_public_keys"` // A list of GPG public keys.
}

type MetaFile struct {
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

func createMetaFiles(provider provider.Provider, file provider.RepositoryFile) error {
	for _, d := range file.Versions {
		version := d.Version
		for _, dp := range d.Targets {
			// TODO refactor
			platform := Platform{
				OS:   dp.OS,
				Arch: dp.Arch,
			}
			filePath := getMetaFilePath(provider, version, platform)

			// TODO Keys
			fileContent := MetaFile{
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

			err := files.WriteToFile(filePath, fileContent)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getMetaFilePath(provider provider.Provider, version string, platform Platform) string {
	return fmt.Sprintf("dist/v1/providers/%s/%s/%s/download/%s/%s", provider.Namespace, provider.ProviderName, internal.TrimTagPrefix(version), platform.OS, platform.Arch)
}
