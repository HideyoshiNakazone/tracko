package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
)

func Test_ExecuteConfigRepoRemove(t *testing.T) {
	// Prepare config
	expectedConfig, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor("Test User", []string{"test@example.com"}).
		WithTargetRepo("test/repo").
		WithTrackedRepos([]string{
			"/tmp/repo1",
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
			"config", "repo", "remove", "/tmp/repo1",
		},
	)

	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	cfg, err := config_handler.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if len(cfg.TrackedRepos()) != 0 {
		t.Errorf("Expected 0 tracked repos, got %d", len(cfg.TrackedRepos()))
	}
}

func Test_ExecuteConfigRepoRemove_IfNotExists(t *testing.T) {
	// Prepare config
	expectedConfig, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor("Test User", []string{"test@example.com"}).
		WithTargetRepo("test/repo").
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
			"config", "repo", "remove", "/tmp/repo1",
		},
	)

	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err == nil {
		t.Fatalf("Command execution succeeded unexpectedly")
	}
}

func Test_ExecuteConfigRepoRemove_InvalidPath(t *testing.T) {
	// Prepare config
	expectedConfig, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor("Test User", []string{"test@example.com"}).
		WithTargetRepo("test/repo").
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
			"config", "repo", "add", "/invalid/path",
		},
	)

	var outputBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outputBuf)

	if err := cmd.RootCmd.Execute(); err == nil {
		t.Fatalf("Command execution succeeded unexpectedly")
	}
}
