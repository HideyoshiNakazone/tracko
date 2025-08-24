package cmd

import (
	"github.com/HideyoshiNakazone/tracko/cmd/flags"
	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/HideyoshiNakazone/tracko/lib/model"
	"github.com/HideyoshiNakazone/tracko/lib/utils"
	"github.com/spf13/cobra"
)

var (
	dbPath              string
	trackedAuthorName   string
	trackedAuthorEmails []string
	targetRepo          string
)

var ConfigInitCmd = &cobra.Command{
	Use:               "init",
	Long:              `Initialize the configuration of the Tracko CLI.`,
	PersistentPreRunE: runConfigInit,
	RunE:              afterConfigInit,
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	err := config.PrepareConfig(flags.GetConfigPath())

	if err == nil {
		cmd.Println("Configuration already initialized.")
		return nil
	}

	cmd.Println("Initializing configuration...")
	var cfgBuilder = model.NewConfigBuilder()

	if dbPath == "" {
		cmd.Println("Using default database path: ", model.DefaultDBPath)
		dbPath = model.DefaultDBPath
	}
	cfgBuilder.WithDBPath(dbPath)

	if trackedAuthorName == "" {
		utils.ReadStringInto("Git author name: ", &trackedAuthorName)
	}
	cfgBuilder.WithTrackedAuthorName(trackedAuthorName)

	if len(trackedAuthorEmails) == 0 {
		utils.ReadStringSliceInto("Git author emails (comma-separated): ", &trackedAuthorEmails)
	}
	cfgBuilder.WithTrackedAuthorEmails(trackedAuthorEmails)

	if targetRepo == "" {
		utils.ReadStringInto("Target repository (owner/repo): ", &targetRepo)
	}
	cfgBuilder.WithTargetRepo(targetRepo)

	if cfg, err := cfgBuilder.Build(); err == nil {
		config.SetConfig(cfg)
		return nil
	}
	return err
}

func afterConfigInit(cmd *cobra.Command, args []string) error {
	_, err := config.GetConfig()
	if err != nil {
		cmd.Println("There was an error initializing the configuration, please remove the config file and try again.")
		return err
	}
	cmd.Println("Congratulations! The configuration has been initialized.")
	return nil
}

func init() {
	// Initialize flags and configuration for the command
	ConfigInitCmd.Flags().StringVar(&dbPath, "db-path", "", "Path to the database file")
	ConfigInitCmd.Flags().StringVar(&trackedAuthorName, "author-name", "", "Name of the author to track")
	ConfigInitCmd.Flags().StringSliceVar(&trackedAuthorEmails, "author-emails", []string{}, "Emails of the authors to track")
	ConfigInitCmd.Flags().StringVar(&targetRepo, "target-repo", "", "Target repository (owner/repo)")
}
