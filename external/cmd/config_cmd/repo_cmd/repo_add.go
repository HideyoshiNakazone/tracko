package repo_cmd

import (
	"errors"
	"path/filepath"

	"github.com/HideyoshiNakazone/tracko/lib/config_handler"
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
	cfg, err := config_handler.GetConfig()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	repoPath, err := filepath.Abs(args[0])
	if err != nil {
		return errors.New("invalid repository path")
	}

	isValid := repo.IsGitRepository(repoPath)
	if !isValid {
		return errors.New("invalid git repository")
	}

	newCfg, err := cfg.AppendTrackedRepo(repoPath)
	if err != nil {
		return err
	}

	return config_handler.SetConfig(newCfg)
}
