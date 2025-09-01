package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config_handler"
	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
)


func Test_ExecuteImport(t *testing.T) {
	numberOfCommits := 100

	testAuthor := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath, cleanup, err := repo.PrepareTestRepo(testAuthor, numberOfCommits)
	if err != nil {
		t.Fatalf("Failed to prepare test repo: %v", err)
	}
	defer (*cleanup)()

	// Prepare config
	expectedConfig, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor(testAuthor.Name(), testAuthor.Emails()).
		WithTargetRepo("test/repo").
		WithTrackedRepos([]string{*repoPath}).
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
			"import",
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
