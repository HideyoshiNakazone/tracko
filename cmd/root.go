package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/cmd/config_cmd"
	"github.com/HideyoshiNakazone/tracko/cmd/flags"
	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/HideyoshiNakazone/tracko/lib/internal_errors"
)

var rootCmd = &cobra.Command{
	Use:   "tracko",
	Short: "CLI for importing Git commit history",
	Long: `Tracko is a command-line tool for importing Git commit history into a database. It allows you to extract commit metadata from private repositories and store it in a database or consolidate it into a single repository.
This tool is useful for developers working at companies that do not use GitHub and want to keep track of their work history.`,
	PersistentPreRun: initConfig,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(ImportCmd)
	cmd.AddCommand(ExportCmd)
	cmd.AddCommand(config_cmd.ConfigCmd)
}

func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&flags.ConfigPath, "config", "", "Path to the config file")
}

func initConfig(cmd *cobra.Command, args []string) {
	err := config.PrepareConfig(flags.GetConfigPath())
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, internal_errors.ErrConfigNotInitialized):
		cmd.Println("Configuration is not initialized. Please run 'tracko config init' to initialize.")
	default:
		cmd.Println("Error initializing configuration:", err)
		// Handle other errors
	}

	os.Exit(1)
}

func init() {
	addSubCommands(rootCmd)
	addFlags(rootCmd)
}
