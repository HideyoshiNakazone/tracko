package config_cmd

import (
	"github.com/spf13/cobra"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)

var ConfigCmd = &cobra.Command{
	Use:  "config",
	Long: `Manage the configuration of the Tracko CLI.`,
	Run:  runConfig,
}

func runConfig(cmd *cobra.Command, args []string) {
	// TODO: Implement import functionality
	cfg, err := config.GetConfig()

	if err != nil {
		cmd.Print("No valid config found.")
		return
	}

	cmd.Print("Importing Git commit history...")
	cmd.Println(cfg)
}

func addSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(ConfigInitCmd)
}

func init() {
	addSubCommands(ConfigCmd)
}
