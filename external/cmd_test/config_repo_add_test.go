package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config"
)


func Test_ExecuteConfigRepoAdd(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TrackedRepos: []string{},
		TargetRepo: "test/repo",
	})

	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "repo", "add", "../..",
		},
	)
	
	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if len(cfg.TrackedRepos) != 1 {
		t.Errorf("Expected 1 tracked repo, got %d", len(cfg.TrackedRepos))
	}
}

func Test_ExecuteConfigRepoAdd_IfAlreadyAdded(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TrackedRepos: []string{},
		TargetRepo: "test/repo",
	})

	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "repo", "add", "../..",
		},
	)
	
	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	if err := cmd.RootCmd.Execute(); err == nil {
		t.Fatalf("Command execution succeeded unexpectedly")
	}
}


func Test_ExecuteConfigRepoAdd_InvalidPath(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TrackedRepos: []string{},
		TargetRepo: "test/repo",
	})

	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "repo", "add", "/invalid/path",
		},
	)

	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err == nil {
		t.Fatalf("Command execution succeeded unexpectedly")
	}
}
