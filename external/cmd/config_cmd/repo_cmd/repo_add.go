package repo_cmd

import (
	"errors"

	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"github.com/spf13/cobra"
)

var RepoAddCmd = &cobra.Command{
	Use:   "add [REPO]",
	Short: "Add a repository to the tracked list",
	Args:  cobra.ExactArgs(1),
	RunE:  runRepoAdd,
}

func runRepoAdd(cmd *cobra.Command, args []string) error {
	repoPath := args[0]

	repoPath, ok := repo.IsGitRepository(&repoPath)
	if !ok {
		return errors.New("Invalid git repository")
	}

	cmd.Println("Adding repository:", repoPath)

	return nil
}
