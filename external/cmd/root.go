package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/external/cmd/config_cmd"
	"github.com/HideyoshiNakazone/tracko/external/flags"
	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	"github.com/HideyoshiNakazone/tracko/lib/internal_errors"
)

var RootCmd = &cobra.Command{
	Use:   "tracko",
	Short: "CLI for importing Git commit history",
	Long: `Tracko is a command-line tool for importing Git commit history into a database. It allows you to extract commit metadata from private repositories and store it in a database or consolidate it into a single repository.
This tool is useful for developers working at companies that do not use GitHub and want to keep track of their work history.`,
	PersistentPreRunE: initConfig,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command, args []string) error {
	err := config_handler.PrepareConfig(flags.GetConfigPath())
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, internal_errors.ErrConfigNotInitialized):
		return fmt.Errorf("configuration is not initialized, please run 'tracko config init' to initialize")
	default:
		return fmt.Errorf("error initializing configuration: %w", err)
		// Handle other errors
	}
}

func init() {
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(config_cmd.ConfigCmd)

	RootCmd.PersistentFlags().StringVar(&flags.ConfigPath, "config", "", "Path to the config file")

}
