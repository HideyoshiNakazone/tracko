package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config"
)

func Test_RunConfigSet(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TargetRepo: "test/repo",
	})
	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	expectedDBPath := "/tmp/test1.db"

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "set", "db_path", expectedDBPath,
		},
	)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	actualDBPath, err := config.GetConfigAttr[string]("db_path")
	if err != nil {
		t.Fatalf("Failed to get config attribute: %v", err)
	}

	if actualDBPath != expectedDBPath {
		t.Errorf("Expected db_path to be %q, but got %q", expectedDBPath, actualDBPath)
	}
}

func Test_RunConfigSet_InvalidKey(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TargetRepo: "test/repo",
	})
	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "set", "invalid_key", "some_value",
		},
	)

	err = cmd.RootCmd.Execute()
	if err == nil {
		t.Fatalf("Expected command to fail with invalid key, but it succeeded")
	}

	expectedErrorMsg := "field \"invalid_key\" does not exist"
	if !bytes.Contains([]byte(err.Error()), []byte(expectedErrorMsg)) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedErrorMsg, err.Error())
	}
}

func Test_RunConfigSet_RestrictedKey(t *testing.T) {
	tempFile, tempCleanup, err := config.PrepareTestConfig(&config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TargetRepo: "test/repo",
	})
	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()

	cmd.RootCmd.SetArgs(
		[]string{
			"--config", tempFile.Name(),
			"config", "set", "version", "restricted_value",
		},
	)

	err = cmd.RootCmd.Execute()
	if err == nil {
		t.Fatalf("Expected command to fail with restricted key, but it succeeded")
	}

	expectedErrorMsg := "field \"version\" is restricted and cannot be modified"
	if !bytes.Contains([]byte(err.Error()), []byte(expectedErrorMsg)) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedErrorMsg, err.Error())
	}

	actualVersion, err := config.GetConfigAttr[string]("version")
	if err != nil {
		t.Fatalf("Failed to get config attribute: %v", err)
	}

	if actualVersion != "v1" {
		t.Errorf("Expected version to be %q, but got %q", "v1", actualVersion)
	}
}
