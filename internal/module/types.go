package module

type Version struct {
	Version string `json:"version"`
}

type RepositoryFile struct {
	Versions []Version `json:"versions"`
}

type Module struct {
	Namespace string
	Name      string
	System    string
}
