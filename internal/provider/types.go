package provider

type RepositoryFile struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version             string   `json:"version"`
	Protocols           []string `json:"protocols"`
	SHASumsURL          string   `json:"shasums_url"`
	SHASumsSignatureURL string   `json:"shasums_signature_url"`
	Repository          string   `json:"repository"`
	Targets             []Target `json:"targets"`
}

type Target struct {
	OS          string `json:"os"`           // The operating system for which the provider is built.
	Arch        string `json:"arch"`         // The architecture for which the provider is built.
	Filename    string `json:"filename"`     // The filename of the provider binary.
	DownloadURL string `json:"download_url"` // The direct URL to download the provider binary.
	SHASum      string `json:"shasum"`       // The SHA checksum of the provider binary.
}

type Provider struct {
	ProviderName string
	Namespace    string
}
