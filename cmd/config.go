package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)

var ConfigCmd = &cobra.Command{
	Use:  "config",
	Long: `Manage the configuration of the Tracko CLI.`,
	RunE:  runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// TODO: Implement import functionality
	cfg, err := config.GetConfig()

	if err != nil {
		return fmt.Errorf("no valid config found: %w", err)
	}

		// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Field", "Value"})

	table.Append([]string{"Version", cfg.Version})
	table.Append([]string{"DBPath", cfg.DBPath})
	table.Append([]string{"Author Name", cfg.TrackedAuthor.Name})
	table.Append([]string{"Author Emails", fmt.Sprintf("%v", cfg.TrackedAuthor.Emails)})
	table.Append([]string{"Target Repo", cfg.TargetRepo})
	table.Append([]string{"Tracked Repos", fmt.Sprintf("%v", cfg.TrackedRepos)})

	table.Render()
	
	return nil
}

func init() {
	ConfigCmd.AddCommand(ConfigInitCmd)
	ConfigCmd.AddCommand(ConfigSetCmd)
	ConfigCmd.AddCommand(ConfigGetCmd)
	ConfigCmd.AddCommand(ConfigRepo)
}
