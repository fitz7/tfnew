package releases

import (
	"encoding/json"
	"fmt"
)

type ProviderData struct {
	Name    string
	Version string
}
type TerraformData struct {
	Name    string
	Version string
}

func GetLatestProviderRelease(provider string) (*ProviderData, error) {
	providerResponse, err := makeGetRequest(fmt.Sprintf("https://registry.terraform.io/v1/providers/%s", provider))
	if err != nil {
		return nil, err
	}
	var providerData ProviderData
	err = json.Unmarshal(providerResponse, &providerData)
	if err != nil {
		return nil, err
	}

	return &providerData, nil

}

func GetLatestTerraformRelease() (*TerraformData, error) {
	terraformResponse, err := makeGetRequest("https://api.releases.hashicorp.com/v1/releases/terraform/latest")
	if err != nil {
		return nil, err
	}
	var terraformData TerraformData
	err = json.Unmarshal(terraformResponse, &terraformData)
	if err != nil {
		return nil, err
	}
	return &terraformData, nil
}
