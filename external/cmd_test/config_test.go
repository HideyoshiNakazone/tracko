package cmd

import (
	"bytes"
	"testing"

	"github.com/HideyoshiNakazone/tracko/external/cmd"
	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
)

func Test_ExecuteConfigCommand(t *testing.T) {
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

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Valid config get",
			args:    []string{"--config", tempFile.Name(), "config"},
			wantErr: false,
		},
		{
			name:    "Invalid config get",
			args:    []string{"--config", "invalid.yaml", "config"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.RootCmd.SetArgs(tt.args)

			cmd_output := new(bytes.Buffer)
			cmd.RootCmd.SetOut(cmd_output)
			cmd.RootCmd.SetErr(cmd_output)

			err := cmd.RootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
