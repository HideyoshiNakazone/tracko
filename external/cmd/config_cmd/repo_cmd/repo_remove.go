package repo_cmd

import (
	"errors"
	"path/filepath"

	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	"github.com/spf13/cobra"
)

var RepoRemoveCmd = &cobra.Command{
	Use:  "remove",
	Long: `Remove a tracked repository from the configuration of the Tracko CLI.`,
	RunE: runRepoRemove,
}

func runRepoRemove(cmd *cobra.Command, args []string) error {
	cfg, err := config_handler.GetConfig()
	if err != nil {
		return err
	}

	repoPath, err := filepath.Abs(args[0])
	if err != nil {
		return errors.New("invalid repository path")
	}

	newCfg, err := cfg.RemoveTrackedRepo(repoPath)
	if err != nil {
		return err
	}

	return config_handler.SetConfig(newCfg)
}
