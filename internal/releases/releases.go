package releases

import (
	"fmt"
)

func GetLatestProviderRelease(provider string) (map[string]interface{}, error) {
	return makeGetRequest(fmt.Sprintf("https://registry.terraform.io/v1/providers/%s", provider))
}

func GetLatestTerraformRelease() (map[string]interface{}, error) {
	return makeGetRequest("https://api.releases.hashicorp.com/v1/releases/terraform/latest")
}
