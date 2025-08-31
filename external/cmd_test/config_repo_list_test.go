package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config_handler"
	"github.com/HideyoshiNakazone/tracko/lib/config_model"
)


func Test_ExecuteConfigRepoList(t *testing.T) {
	expectedConfig := &config_model.ConfigModel{
		Version: config_model.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config_model.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TrackedRepos: []string{
			"/path/to/your/repo1",
			"/path/to/your/repo2",
		},
		TargetRepo: "test/repo",
	}

	tempFile, tempCleanup, err := config_handler.PrepareTestConfig(expectedConfig)

	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "repo", "list",
		},
	)
	
	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}
}
