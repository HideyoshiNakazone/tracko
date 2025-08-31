package cmd

import (
	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/lib/config_handler"
)

var ImportCmd = &cobra.Command{
	Use:  "import",
	Long: `Import Git commit history from a repository.`,
	Run:  runImport,
}

func runImport(cmd *cobra.Command, args []string) {
	// TODO: Implement import functionality
	cfg, err := config_handler.GetConfig()

	if err != nil {
		cmd.Print("No valid config found.")
		return
	}

	cmd.Print("Importing Git commit history...")
	cmd.Println(cfg)
}
