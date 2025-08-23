package cmd

import (
	"bytes"
	"os"
	"testing"
)

func Test_ExecuteConfigInit(t *testing.T) {
	tempFile, err := os.CreateTemp("", "tracko_test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	cmd_output := new(bytes.Buffer)

	RootCmd.SetOut(cmd_output)
	RootCmd.SetErr(cmd_output)
	RootCmd.SetArgs(
		[]string{
			"config", "init",
			"--config", tempFile.Name(),
			"--db-path", "/tmp/test.db",
			"--author-name", "Test User",
			"--author-emails", "test@example.com",
		},
	)

	if err := RootCmd.Execute(); err != nil {
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
