package config_handler

import (
	"os"
	"testing"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/spf13/viper"
)

func TestSetAndGetConfig(t *testing.T) {
	// Setup: create a temp config file
	tempFile, err := os.CreateTemp("", "tracko_test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Prepare config
	cfg, err := config_model.NewConfigBuilder().
		WithDBPath("/tmp/test.db").
		WithTrackedAuthor("Test User", []string{"test@example.com"}).
		WithTargetRepo("repo1").
		WithTrackedRepos([]string{"repo1", "repo2"}).
		Build()

	if err != nil {
		t.Fatalf("Failed to build config: %v", err)
	}

	// Set config file for viper
	viper.SetConfigFile(tempFile.Name())
	viper.SetConfigType("yaml")

	// Test SetConfig
	err = SetConfig(cfg)
	if err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	// Test GetConfig
	got, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if got.Version() != cfg.Version() || got.DBPath() != cfg.DBPath() || got.TrackedAuthor().Name() != cfg.TrackedAuthor().Name() {
		t.Errorf("Config values do not match. Got: %+v, Want: %+v", got, cfg)
	}
	if len(got.TrackedRepos()) != len(cfg.TrackedRepos()) {
		t.Errorf("TrackedRepos length mismatch. Got: %d, Want: %d", len(got.TrackedRepos()), len(cfg.TrackedRepos()))
	}
}

func TestPrepareConfigWithInvalidFile(t *testing.T) {
	err := PrepareConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent config file, got nil")
	}
}

func Test_SetConfigAttr(t *testing.T) {
	// Test SetConfigAttr
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Restricted case", args{"version", "value1"}, true},
		{"Valid case", args{"db_path", "value1"}, false},
		{"Invalid case", args{"invalid_field", "value2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config_model.NewConfigBuilder().
				WithDBPath("/tmp/test.db").
				WithTrackedAuthor("Test User", []string{"test@example.com"}).
				WithTargetRepo("repo1").
				WithTrackedRepos([]string{"repo1", "repo2"}).
				Build()

			if err != nil {
				t.Fatalf("Failed to build config: %v", err)
			}

			_, cleanup, err := PrepareTestConfig(cfg)
			defer (*cleanup)()

			if err != nil {
				t.Fatalf("Failed to prepare test config: %v", err)
			}

			if err := SetConfigAttr(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetConfigAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
