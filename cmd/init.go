package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Version: rootCmd.Version,
	Short:   "creates a config file in the repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.SafeWriteConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
