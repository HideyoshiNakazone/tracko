package cmd

import (
	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)


var ExportCmd = &cobra.Command{
	Use:   "export",
	Long:  `Export Git commit history to a repository.`,
	Run: runExport,
}


func runExport(cmd *cobra.Command, args []string) {
	// TODO: Implement import functionality
	cfg, err := config.GetConfig();

	if err != nil {
		cmd.Print("No valid config found.")
		return
	}

	cmd.Print("Importing Git commit history...")
	cmd.Println(cfg)
}
