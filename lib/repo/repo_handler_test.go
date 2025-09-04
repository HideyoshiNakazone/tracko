package repo

import (
	"testing"

	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
)

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

func Test_NewTrackedRepo(t *testing.T) {
	author := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	tests := []struct {
		name string
		args *struct {
			author *config_model.ConfigAuthorModel
			path   string
		}
		validator func(*TrackedRepo) bool
		wantErr   bool
	}{
		{
			name: "valid repo",
			args: &struct {
				author *config_model.ConfigAuthorModel
				path   string
			}{
				author: author,
				path:   "../..",
			},
			validator: func(tracked *TrackedRepo) bool {
				return tracked != nil && tracked.repo != nil && tracked.author != nil
			},
			wantErr: false,
		},
		{
			name: "invalid repo",
			args: &struct {
				author *config_model.ConfigAuthorModel
				path   string
			}{
				author: author,
				path:   "/invalid/path",
			},
			validator: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracked, err := NewTrackedRepo(tt.args.path, tt.args.author)
			if (tt.wantErr && err == nil) || (!tt.wantErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && tracked == nil {
				t.Errorf("expected non-nil repo, got nil")
			}

			if tt.validator != nil && !tt.validator(tracked) {
				t.Errorf("failed validation for repo: %v", tracked)
			}
		})
	}
}

func Test_ListRepositoryHistory_With_ForEach(t *testing.T) {
	numberOfCommits := 100

	testAuthor := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath, cleanup, err := PrepareTestRepo(testAuthor, numberOfCommits)
	if err != nil {
		t.Fatalf("Failed to prepare test repo: %v", err)
	}
	defer (*cleanup)()

	tests := []struct {
		name          string
		author        *config_model.ConfigAuthorModel
		params        *ListRepositoryHistoryParams
		expectedCount int
	}{
		{
			name:          "Test User",
			author:        testAuthor,
			params:        &ListRepositoryHistoryParams{},
			expectedCount: numberOfCommits,
		},
		{
			name: "Other User",
			author: config_model.AuthorDTO{
				Name: "Invalid User",
				Emails: []string{
					"invalid_user@example.com",
				},
			}.ToModel(),
			params:        &ListRepositoryHistoryParams{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracked, err := NewTrackedRepo(*repoPath, tt.author)
			if err != nil {
				t.Fatalf("Failed to create tracked repo: %v", err)
			}

			iter, err := tracked.ListRepositoryHistory(tt.params)
			if err != nil {
				t.Fatalf("Failed to get commit iterator: %v", err)
			}
			defer iter.Close()

			commitCount := 0
			iter.ForEach(func(meta *GitCommitMeta) error {
				commitCount++
				return nil
			})

			if commitCount != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d", tt.expectedCount, commitCount)
			}
		})
	}
}

func Test_ListRepositoryHistory_With_Next(t *testing.T) {
	numberOfCommits := 100

	testAuthor := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath, cleanup, err := PrepareTestRepo(testAuthor, numberOfCommits)
	if err != nil {
		t.Fatalf("Failed to prepare test repo: %v", err)
	}
	defer (*cleanup)()

	tests := []struct {
		name          string
		author        *config_model.ConfigAuthorModel
		params        *ListRepositoryHistoryParams
		expectedCount int
	}{
		{
			name:          "Test User",
			author:        testAuthor,
			params:        &ListRepositoryHistoryParams{},
			expectedCount: numberOfCommits,
		},
		{
			name: "Other User",
			author: config_model.AuthorDTO{
				Name: "Invalid User",
				Emails: []string{
					"invalid_user@example.com",
				},
			}.ToModel(),
			params:        &ListRepositoryHistoryParams{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracked, err := NewTrackedRepo(*repoPath, tt.author)
			if err != nil {
				t.Fatalf("Failed to create tracked repo: %v", err)
			}

			iter, err := tracked.ListRepositoryHistory(tt.params)
			if err != nil {
				t.Fatalf("Failed to get commit iterator: %v", err)
			}
			defer iter.Close()

			commitCount := 0
			for {
				_, err := iter.Next()
				if err != nil {
					break
				}
				commitCount++
			}

			if commitCount != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d", tt.expectedCount, commitCount)
			}
		})
	}
}
