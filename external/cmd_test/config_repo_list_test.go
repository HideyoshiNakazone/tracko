package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
)

func Test_ExecuteConfigRepoList(t *testing.T) {
	// Prepare config
	expectedConfig, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor("Test User", []string{"test@example.com"}).
		WithTargetRepo("test/repo").
		WithTrackedRepos([]string{
			"/path/to/your/repo1",
			"/path/to/your/repo2",
		}).
		Build()

	if err != nil {
		t.Fatalf("Failed to build expected config: %v", err)
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
