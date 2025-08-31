package config_model

import (
	"reflect"
	"testing"
)


func Test_ConfigBuilder(t *testing.T) {
    tests := []struct {
        name     	string
        builder    	*ConfigModelBuilder
        expected 	*ConfigModel
		wantErr  	bool
    }{
        {
            name:     "valid config",
            builder:  NewConfigBuilder().
						WithDBPath("/tmp/test.db").
						WithTrackedAuthor("Test User", []string{
							"test@example.com",
						}).
						WithTargetRepo("repo1"),
            expected: &ConfigModel{
				version: 	 	CurrentVersion,
				dbPath:        	"/tmp/test.db",
				trackedAuthor: ConfigAuthorModel{
					name:  "Test User",
					emails: []string{"test@example.com"},
				},
				targetRepo:   	"repo1",
				trackedRepos: 	[]string{},
			},
			wantErr: 	false,
        },
        {
            name:    	"invalid config - missing target repo",
            builder: 	NewConfigBuilder().
							WithTrackedAuthor("Test User", []string{
								"test@example.com",
							}),
            expected: 	nil,
            wantErr: 	true,
        },
        {
            name:    	"invalid config - missing tracked author",
            builder: 	NewConfigBuilder().
							WithTargetRepo("repo1"),
            expected: 	nil,
            wantErr: 	true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            actualModel, err := tt.builder.Build()
			if (tt.wantErr && err == nil) || (!tt.wantErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !reflect.DeepEqual(actualModel, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actualModel)
			}
        })
    }
}
