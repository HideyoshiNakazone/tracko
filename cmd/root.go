package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)

var (
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "tracko",
	Short: "CLI for importing Git commit history",
	Long: `Tracko is a command-line tool for importing Git commit history into a database. It allows you to extract commit metadata from private repositories and store it in a database or consolidate it into a single repository.
This tool is useful for developers working at companies that do not use GitHub and want to keep track of their work history.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(ImportCmd)
	cmd.AddCommand(ExportCmd)
}

func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the config file")
}

func initConfig() {
	var err error;

	if configPath == "" {
		err = config.InitializeConfig()
	} else {
		err = config.InitializeConfigFromFile(configPath)
	}

	if err != nil {
		panic(err)
	}
}

func init() {
	addSubCommands(rootCmd)
	addFlags(rootCmd)
	cobra.OnInitialize(initConfig)
}
