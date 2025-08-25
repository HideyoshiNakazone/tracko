package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/lib/config"
)


func Test_ExecuteConfigCommand(t *testing.T) {
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

	tests := []struct {
		name       string
		args      []string
		wantErr   bool
	}{
		{
			name:     "Valid config get",
			args:    []string{"--config", tempFile.Name(), "config"},
			wantErr: false,
		},
		{
			name:     "Invalid config get",
			args:    []string{"--config", "invalid.yaml", "config"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RootCmd.SetArgs(tt.args)

			cmd_output := new(bytes.Buffer)
			RootCmd.SetOut(cmd_output)
			RootCmd.SetErr(cmd_output)

			err := RootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
