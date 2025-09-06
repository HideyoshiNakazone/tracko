package config_model

import (
	"reflect"
	"testing"

	"github.com/HideyoshiNakazone/tracko/lib/internal_errors"
)

func Test_ConfigModelBuilder_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  *ConfigModel
		want    *ConfigModel
		wantErr error
	}{
		{
			name: "valid config",
			config: &ConfigModel{
				version:       "v1",
				dbPath:        "$HOME/.config/tracko.db",
				trackedAuthor: ConfigAuthorModel{name: "test", emails: []string{"test@example.com"}},
				targetRepo:    "test/repo",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			want: &ConfigModel{
				version:       "v1",
				dbPath:        "$HOME/.config/tracko.db",
				trackedAuthor: ConfigAuthorModel{name: "test", emails: []string{"test@example.com"}},
				targetRepo:    "test/repo",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			wantErr: nil,
		},
		{
			name: "invalid config - missing version",
			config: &ConfigModel{
				version:       "",
				dbPath:        "$HOME/.config/tracko.db",
				trackedAuthor: ConfigAuthorModel{name: "test", emails: []string{"test@example.com"}},
				targetRepo:    "test/repo",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			want:    nil,
			wantErr: internal_errors.ErrInvalidConfig,
		},
		{
			name: "invalid config - missing db path",
			config: &ConfigModel{
				version:       "v1",
				dbPath:        "",
				trackedAuthor: ConfigAuthorModel{name: "test", emails: []string{"test@example.com"}},
				targetRepo:    "test/repo",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			want:    nil,
			wantErr: internal_errors.ErrInvalidConfig,
		},
		{
			name: "invalid config - missing tracked author",
			config: &ConfigModel{
				version:       "v1",
				dbPath:        "$HOME/.config/tracko.db",
				trackedAuthor: ConfigAuthorModel{name: "", emails: []string{"test@example.com"}},
				targetRepo:    "test/repo",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			want:    nil,
			wantErr: internal_errors.ErrInvalidConfig,
		},
		{
			name: "invalid config - missing target repo",
			config: &ConfigModel{
				version:       "v1",
				dbPath:        "$HOME/.config/tracko.db",
				trackedAuthor: ConfigAuthorModel{name: "test", emails: []string{"test@example.com"}},
				targetRepo:    "",
				trackedRepos:  []string{"repo1", "repo2"},
			},
			want:    nil,
			wantErr: internal_errors.ErrInvalidConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &ConfigModelBuilder{config: tt.config}
			got, err := builder.Build()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %+v, want %+v", got, tt.want)
			}

			if err != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
