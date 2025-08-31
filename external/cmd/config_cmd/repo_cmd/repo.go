package repo_cmd

import (
	"github.com/spf13/cobra"
)

var RepoCmd = &cobra.Command{
	Use:  "repo",
	Long: `Manage the tracked repositories in the configuration of the Tracko CLI.`,
}


func init() {
	RepoCmd.AddCommand(RepoListCmd)
	RepoCmd.AddCommand(RepoAddCmd)
	RepoCmd.AddCommand(RepoRemoveCmd)
}
