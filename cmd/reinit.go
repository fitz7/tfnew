package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var reinitCmd = &cobra.Command{
	Use:     "reinit",
	Version: rootCmd.Version,
	Short:   "reinitialise your config file to the defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.WriteConfig()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reinitCmd)
}
