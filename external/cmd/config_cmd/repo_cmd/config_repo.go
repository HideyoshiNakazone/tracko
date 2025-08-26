package repo_cmd

import (
	"errors"

	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/spf13/cobra"
)

var ConfigRepo = &cobra.Command{
	Use:  "repo",
	Long: `Manage the tracked repositories in the configuration of the Tracko CLI.`,
	RunE: runConfigRepo,
}

func runConfigRepo(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		cmd.Println("Invalid number of arguments. Usage: config repo <subcommand>")
		return errors.New("invalid number of arguments")
	}

	switch args[0] {
		case "list":
			return runConfigRepoList(cmd, args[1:])
		// case "add":
		// 	return runConfigRepoAdd(cmd, args[1:])
		// case "remove":
		// 	return runConfigRepoRemove(cmd, args[1:])
		default:
			cmd.Println("Unknown subcommand:", args[0])
			return errors.New("unknown subcommand")
	}
}


func runConfigRepoList(cmd *cobra.Command, args []string) error {
	repos, err := config.GetConfigAttr[[]interface{}]("tracked_repos")
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		cmd.Println("No tracked repositories found.")
		return nil
	}

	cmd.Println("Tracked repositories:")
	for _, repo := range repos {
		cmd.Println("-", repo)
	}
	return nil
}
