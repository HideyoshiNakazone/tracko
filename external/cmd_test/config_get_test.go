package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	"github.com/HideyoshiNakazone/tracko/lib/config"
)


func Test_ExecuteConfigGet(t *testing.T) {
	expectedConfig := &config.ConfigModel{
		Version: config.CurrentVersion,
		DBPath:  "/tmp/test.db",
		TrackedAuthor: config.ConfigAuthorModel{
			Name: "Test User",
			Emails: []string{
				"test@example.com",
			},
		},
		TargetRepo: "test/repo",
	}

	tempFile, tempCleanup, err := config.PrepareTestConfig(expectedConfig)
	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer (*tempCleanup)()
	
	tests := []struct {
		name           string
		key           string
		expectedValue any
		wantErr        bool
	}{
		{
			name: "Get DB Path",
			key:  "db_path",
			expectedValue: expectedConfig.DBPath,
			wantErr:        false,
		},
		{
			name: "Get Version",
			key:  "version",
			expectedValue: expectedConfig.Version,
			wantErr:        false,
		},
		{
			name: "Get Tracked Author Name",
			key:  "author.name",
			expectedValue: expectedConfig.TrackedAuthor.Name,
			wantErr:        false,
		},
		{
			name: "Get Tracked Author Emails",
			key:  "author.emails",
			expectedValue: expectedConfig.TrackedAuthor.Emails,
			wantErr:        false,
		},
		{
			name: "Get Target Repo",
			key:  "target_repo",
			expectedValue: expectedConfig.TargetRepo,
			wantErr:        false,
		},
		{
			name: "Get Tracked Repos",
			key:  "tracked_repos",
			expectedValue: expectedConfig.TrackedRepos,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedOutput := fmt.Sprintf("%s => %v\n", tt.key, tt.expectedValue)

			cmd.RootCmd.SetArgs(
				[]string{
					"--config", tempFile.Name(),
					"config", "get", tt.key,
				},
			)

			cmd_output := new(bytes.Buffer)

			cmd.RootCmd.SetOut(cmd_output)
			cmd.RootCmd.SetErr(cmd_output)

			err := cmd.RootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if cmd_output.String() != expectedOutput {
				t.Errorf("Expected output %q, but got %q", expectedOutput, cmd_output.String())
			}
		})
	}
}
