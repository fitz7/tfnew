package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fitz7/tfnew/internal/fsutils"
)

var (
	Version = "dev"
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tfnew",
	Version: Version,
	Short:   "Creates terraform modules with a standard structure",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .tfnew.yaml)")

	rootCmd.PersistentFlags().String("backend", "local", "The backend to use in your new module (default is 'local'")
	_ = viper.BindPFlag("backend.type", rootCmd.PersistentFlags().Lookup("backend"))

	// local backend flags
	rootCmd.PersistentFlags().String("backend_local_path", "./terraform.tfstate", "the path to use for your local backend")

	// gcs backend flags
	rootCmd.PersistentFlags().String("backend_gcs_bucket", "", "the bucket to use for your gcs backend")
	rootCmd.PersistentFlags().String("backend_gcs_prefix", "", "the root prefix to use for your gcs backend (default is ''")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		// for the init command to create the file
		viper.AddConfigPath(filepath.Dir(cfgFile))
		viper.SetConfigName(filepath.Base(cfgFile[:len(cfgFile)-len(filepath.Ext(cfgFile))]))
		viper.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgFile), "."))
	} else {
		projectRoot := fsutils.FindProjectRootDir()

		viper.AddConfigPath(projectRoot)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tfnew")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// the backend dependant bindings are done here so that the file created by init is not polluted with backend types that are not relevant
	backendBindings()
}

func backendBindings() {
	backendType := viper.GetString("backend.type")

	switch backendType {
	case "local":
		_ = viper.BindPFlag("backend.local.path", rootCmd.PersistentFlags().Lookup("backend_local_path"))
	case "gcs":
		_ = viper.BindPFlag("backend.gcs.bucket", rootCmd.PersistentFlags().Lookup("backend_gcs_bucket"))
		_ = viper.BindPFlag("backend.gcs.prefix", rootCmd.PersistentFlags().Lookup("backend_gcs_root_prefix"))
	case "s3", "remote", "azurerm", "consul", "cos", "http", "kubernetes", "oss", "pg", "cloud":
		log.Fatalf("backend %s not implemented", backendType)
	}
}
