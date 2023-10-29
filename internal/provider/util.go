package provider

import "fmt"

func GetRepositoryName(provider Provider) string {
	return fmt.Sprintf("terraform-provider-%s", provider.ProviderName) // TODO hashi?
}

func GetRepositoryUrl(provider Provider) string {
	return fmt.Sprintf("https://github.com/%s/%s", EffectiveNamespace(provider.Namespace), GetRepositoryName(provider))
}

// EffectiveProviderNamespace will map namespaces for providers in situations
// where the author (owner of the namespace) does not release artifacts as
// GitHub Releases.
func EffectiveNamespace(namespace string) string {
	if namespace == "hashicorp" {
		return "opentofu"
	}

	return namespace
} // TODO make more generic
