package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)


func Test_ExecuteConfigRepoList(t *testing.T) {
	expectedConfig := &config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
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

	tempFile, tempCleanup, err := config.PrepareTestConfig(expectedConfig)

	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "repo", "list",
		},
	)
	
	var outputBuf bytes.Buffer
	RootCmd.SetOut(&outputBuf)

	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}
}
