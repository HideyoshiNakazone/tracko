package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
)

func Test_ExecuteConfigInit(t *testing.T) {
	tempFile, err := os.CreateTemp("", "tracko_test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	cmd_output := new(bytes.Buffer)

	cmd.RootCmd.SetOut(cmd_output)
	cmd.RootCmd.SetErr(cmd_output)
	cmd.RootCmd.SetArgs(
		[]string{
			"config", "init",
			"--config", tempFile.Name(),
			"--db-path", "/tmp/test.db",
			"--author-name", "Test User",
			"--author-emails", "test@example.com",
			"--target-repo", "test/repo",
		},
	)

	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	output := cmd_output.String()
	expectedOutputs := []string{
		"Initializing configuration...",
		"Congratulations! The configuration has been initialized.",
	}

	for _, expected := range expectedOutputs {
		if !bytes.Contains([]byte(output), []byte(expected)) {
			t.Errorf("Expected output to contain '%s', but it did not. Full output: %s", expected, output)
		}
	}
}
