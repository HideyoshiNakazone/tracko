package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config_handler"
	"github.com/HideyoshiNakazone/tracko/lib/config_model"
)


func Test_ExecuteConfigRepoAdd(t *testing.T) {
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
			"config", "repo", "add", "../..",
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

	if len(cfg.TrackedRepos()) != 1 {
		t.Errorf("Expected 1 tracked repo, got %d", len(cfg.TrackedRepos()))
	}
}

func Test_ExecuteConfigRepoAdd_IfAlreadyAdded(t *testing.T) {
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
