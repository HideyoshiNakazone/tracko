package config_cmd

import (
	"errors"

	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/spf13/cobra"
)

var ConfigGetCmd = &cobra.Command{
	Use:  "get",
	Long: `Get an attribute from the configuration of the Tracko CLI.`,
	RunE: runConfigGet,
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		cmd.Println("Invalid number of arguments. Usage: config get <key>")
		return errors.New("invalid number of arguments")
	}

	value, err := config.GetConfigAttr[any](args[0])
	if err != nil {
		return err
	}

	// {key} => {value}
	cmd.Printf("%s => %v\n", args[0], value)
	return nil
}
