package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fitz7/tfnew/internal/tfmodule"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:     "module",
	Version: rootCmd.Version,
	Short:   "Creates a new terraform module",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a single arg containing the path of the new module")
		}
		return validatePath(args[0])
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		requiredProviders, err := cmd.Flags().GetStringSlice("required_providers")
		if err != nil {
			return err
		}

		root, err := cmd.Flags().GetBool("root")
		if err != nil {
			return err
		}

		backend := viper.GetString("backend.type")

		createModuleOptions := tfmodule.CreateModuleOptions{
			Name:              args[0],
			RootModule:        root,
			RequiredProviders: requiredProviders,
			BackendType:       backend,
		}

		fullModulePath, err := tfmodule.CreateModuleDir(createModuleOptions)
		if err != nil {
			return err
		}
		defaultFiles, err := tfmodule.CreateDefaultModuleFiles(fullModulePath)
		if err != nil {
			return fmt.Errorf("error creating moduleName files: %w", err)
		}
		err = tfmodule.PopulateVersionsFile(defaultFiles[tfmodule.VersionsFile], createModuleOptions)
		if err != nil {
			return fmt.Errorf("error populating the versions.tf file: %w", err)
		}

		defer func() {
			for _, file := range defaultFiles {
				_ = file.Close()
			}
		}()
		return nil
	},
}

func validatePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return fmt.Errorf("new module directory already exists: %s", path)
}

func init() {
	rootCmd.AddCommand(moduleCmd)

	moduleCmd.Flags().StringSliceP("required_providers", "p", []string{}, "")

	moduleCmd.Flags().BoolP("root", "r", false, "")
}
