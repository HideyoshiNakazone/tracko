package repo

import "testing"

func Test_IsGitRepository(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Valid Git Repo", "../..", true},
		{"Invalid Git Repo", "/path/to/invalid/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := IsGitRepository(tt.path)
			if ok != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, ok)
			}
		})
	}
}
