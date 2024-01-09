package tfmodule

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"

	"github.com/fitz7/tfnew/internal/fsutils"
	"github.com/fitz7/tfnew/internal/releases"
)

const versionsFile = "versions.tf"

var defaultFilenames = []string{versionsFile, "variables.tf", "outputs.tf", "main.tf"}

type CreateModuleOptions struct {
	Name              string
	RootModule        bool
	RequiredProviders []string
	TerraformVersion  string
	BackendType       string
}

func CreateModule(cmo CreateModuleOptions) error {
	fullPathWithModuleName := fmt.Sprintf("%s/%s", fsutils.FindProjectRootDir(), cmo.Name)

	err := os.Mkdir(fullPathWithModuleName, 0o755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	defaultFiles, err := createDefaultModuleFiles(fullPathWithModuleName)
	if err != nil {
		return fmt.Errorf("error creating moduleName files: %w", err)
	}

	defer func() {
		for _, file := range defaultFiles {
			_ = file.Close()
		}
	}()

	err = populateVersionsFile(defaultFiles[versionsFile], cmo)
	if err != nil {
		return fmt.Errorf("error populating the versions.tf file: %w", err)
	}

	return nil
}

func createDefaultModuleFiles(path string) (map[string]*os.File, error) {
	defaultFiles := make(map[string]*os.File)

	for _, filename := range defaultFilenames {
		newFile, err := os.Create(fmt.Sprintf("%s/%s", path, filename))
		if err != nil {
			return nil, err
		}

		defaultFiles[filename] = newFile
	}

	return defaultFiles, nil
}

func populateVersionsFile(versionsFile *os.File, cmo CreateModuleOptions) error {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	err := addTerraformBlock(body, cmo)
	if err != nil {
		return err
	}

	_, err = f.WriteTo(versionsFile)
	if err != nil {
		return err
	}

	return nil
}

func addTerraformBlock(body *hclwrite.Body, cmo CreateModuleOptions) error {
	terraformBody := body.AppendNewBlock("terraform", nil).Body()

	err := addRequiredVersion(cmo, terraformBody)
	if err != nil {
		return err
	}

	if cmo.RootModule {
		err = addBackendBlock(cmo, terraformBody)
		if err != nil {
			return err
		}
	}

	err = addRequiredProvidersBlock(cmo, terraformBody)
	if err != nil {
		return err
	}

	return nil
}

func addRequiredVersion(cmo CreateModuleOptions, terraformBody *hclwrite.Body) error {
	if cmo.TerraformVersion == "" {
		var err error

		cmo.TerraformVersion, err = getTerraformVersion(cmo.RootModule)
		if err != nil {
			return err
		}
	}

	terraformBody.SetAttributeValue("required_version", cty.StringVal(cmo.TerraformVersion))

	return nil
}

func addBackendBlock(cmo CreateModuleOptions, terraformBody *hclwrite.Body) error {
	backendBody := terraformBody.AppendNewBlock("backend", []string{cmo.BackendType}).Body()

	switch cmo.BackendType {
	case "gcs":
		// get the backendConfig from the config file
		backendConfig := viper.GetStringMap(fmt.Sprintf("backend.%s", cmo.BackendType))
		// fill in any missing keys
		keys := []string{"prefix", "bucket"}
		for _, key := range keys {
			backendConfig[key] = viper.Get(fmt.Sprintf("backend.%s.%s", cmo.BackendType, key))
		}

		// if there is no prefix we just use the name
		if prefix, ok := backendConfig["prefix"]; !ok || prefix == nil || prefix == "" {
			backendConfig["prefix"] = cmo.Name
		} else {
			backendConfig["prefix"] = fmt.Sprintf("%s/%s", prefix, cmo.Name)
		}

		writeTerraformFromAnyMap(backendConfig, backendBody)
	case "local":
		backendBody.SetAttributeValue("path", cty.StringVal("./terraform.tfstate"))
	default:
		return errors.New("backend not implemented")
	}

	return nil
}

func writeTerraformFromAnyMap(anyMap map[string]interface{}, backendBody *hclwrite.Body) {
	for key, value := range anyMap {
		switch typedVal := value.(type) {
		case string:
			backendBody.SetAttributeValue(key, cty.StringVal(typedVal))
		case bool:
			backendBody.SetAttributeValue(key, cty.BoolVal(typedVal))
		//	gcs doesn't have any nested config, but we'll include it here anyway as a POC
		case map[string]interface{}:
			newBody := backendBody.AppendNewBlock(key, []string{}).Body()
			writeTerraformFromAnyMap(typedVal, newBody)
		}
	}
}

func addRequiredProvidersBlock(cmo CreateModuleOptions, body *hclwrite.Body) error {
	if len(cmo.RequiredProviders) == 0 {
		return nil
	}

	requiredProvidersBody := body.AppendNewBlock("required_providers", []string{}).Body()

	for _, provider := range cmo.RequiredProviders {
		latestProviderRelease, err := releases.GetLatestProviderRelease(provider)
		if err != nil {
			return err
		}

		providerName, ok := latestProviderRelease["name"].(string)
		if !ok {
			return fmt.Errorf("could not find name for provider: %s", provider)
		}

		latestProviderVersion, ok := latestProviderRelease["version"].(string)
		if !ok {
			return fmt.Errorf("could not find version for provider: %s", provider)
		}

		minorProviderVersion := truncatePatchVersion(latestProviderVersion)

		requiredProvidersBody.SetAttributeValue(providerName, cty.ObjectVal(map[string]cty.Value{
			"source":  cty.StringVal(provider),
			"version": cty.StringVal(fmt.Sprintf("~> %s", minorProviderVersion)),
		}))
	}

	return nil
}

func getTerraformVersion(rootModule bool) (string, error) {
	// default for non-root modules
	// maybe these should be split
	terraformVersion := ">= 1.0"

	if rootModule {
		// when root module we want to use the latest terraform version, so we fetch and format it
		latestTerraformRelease, err := releases.GetLatestTerraformRelease()
		if err != nil {
			return "", errors.New("failed to fetch terraform version")
		}

		latestTerraformVersion, ok := latestTerraformRelease["version"].(string)
		if !ok {
			return "", err
		}

		minorTerraformVersion := truncatePatchVersion(latestTerraformVersion)

		terraformVersion = fmt.Sprintf("~> %s", minorTerraformVersion)
	}

	return terraformVersion, nil
}

func truncatePatchVersion(version string) string {
	return strings.Join(strings.Split(version, ".")[:2], ".")
}
