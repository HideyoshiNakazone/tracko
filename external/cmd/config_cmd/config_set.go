package config_cmd

import (
	"errors"

	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/spf13/cobra"
)

var ConfigSetCmd = &cobra.Command{
	Use:  "set",
	Long: `Set an attribute in the configuration of the Tracko CLI.`,
	RunE: runConfigSet,
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		cmd.Println("Invalid number of arguments. Usage: config set <key> <value>")
		return errors.New("invalid number of arguments")
	}

	return config.SetConfigAttr(args[0], args[1])
}
