package repo_cmd

import (
	"github.com/HideyoshiNakazone/tracko/lib/config"
	"github.com/spf13/cobra"
)

var RepoListCmd = &cobra.Command{
	Use:  "list",
	Long: `List the tracked repositories in the configuration of the Tracko CLI.`,
	RunE: runRepoList,
}

func runRepoList(cmd *cobra.Command, args []string) error {
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
