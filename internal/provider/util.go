package provider

import "fmt"

func GetRepositoryName(provider Provider) string {
	return fmt.Sprintf("terraform-provider-%s", provider.ProviderName) // TODO hashi?
}
