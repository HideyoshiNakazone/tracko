package cmd

import (
	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:  "export",
	Long: `Export Git commit history to a repository.`,
	Run:  runExport,
}

func runExport(cmd *cobra.Command, args []string) {
	// TODO: Implement import functionality
	cfg, err := config_handler.GetConfig()

	if err != nil {
		cmd.Print("No valid config found.")
		return
	}

	cmd.Print("Importing Git commit history...")
	cmd.Println(cfg)
}
