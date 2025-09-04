package cmd

import (
	"github.com/spf13/cobra"

	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	"github.com/HideyoshiNakazone/tracko/lib/import_handler"
)

var ImportCmd = &cobra.Command{
	Use:  "import",
	Long: `Import Git commit history from a repository.`,
	RunE: runImport,
}

func runImport(cmd *cobra.Command, args []string) error {
	// TODO: Implement import functionality
	cfg, err := config_handler.GetConfig()

	if err != nil {
		cmd.Print("No valid config found.")
		return err
	}

	return import_handler.ImportTrackedRepos(cfg)
}
