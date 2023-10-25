package provider

import "fmt"

func GetRepositoryName(provider Provider) string {
	return fmt.Sprintf("terraform-provider-%s", provider.ProviderName) // TODO hashi?
}

func GetRepositoryUrl(provider Provider) string {
	return fmt.Sprintf("https://github.com/%s/%s", provider.Namespace, GetRepositoryName(provider))
}
