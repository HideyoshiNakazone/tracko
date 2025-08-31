package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
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

	author := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath, cleanup, err := prepareTestRepo(author, numberOfCommits)
	if err != nil {
		t.Fatalf("Failed to prepare test repo: %v", err)
	}
	defer (*cleanup)()

	tests := []struct {
		name          string
		params        *ListRepositoryHistoryParams
		expectedCount int
	}{
		{
			name: "Test User",
			params: &ListRepositoryHistoryParams{
				Author: author,
			},
			expectedCount: numberOfCommits,
		},
		{
			name: "Other User",
			params: &ListRepositoryHistoryParams{
				Author: config_model.AuthorDTO{
					Name: "Invalid User",
					Emails: []string{
						"invalid_user@example.com",
					},
				}.ToModel(),
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracked, err := NewTrackedRepo(*repoPath, author)
			if err != nil {
				t.Fatalf("Failed to create tracked repo: %v", err)
			}

			iter, err := tracked.ListRepositoryHistory(tt.params)
			if err != nil {
				t.Fatalf("Failed to get commit iterator: %v", err)
			}
			defer iter.Close()

			commitCount := 0
			iter.ForEach(func (meta *GitCommitMeta) error {
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

	author := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath, cleanup, err := prepareTestRepo(author, numberOfCommits)
	if err != nil {
		t.Fatalf("Failed to prepare test repo: %v", err)
	}
	defer (*cleanup)()

	tests := []struct {
		name          string
		params        *ListRepositoryHistoryParams
		expectedCount int
	}{
		{
			name: "Test User",
			params: &ListRepositoryHistoryParams{
				Author: author,
			},
			expectedCount: numberOfCommits,
		},
		{
			name: "Other User",
			params: &ListRepositoryHistoryParams{
				Author: config_model.AuthorDTO{
					Name: "Invalid User",
					Emails: []string{
						"invalid_user@example.com",
					},
				}.ToModel(),
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracked, err := NewTrackedRepo(*repoPath, author)
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

func prepareTestRepo(author *config_model.ConfigAuthorModel, numberOfCommits int) (*string, *func(), error) {
	tempDir, err := ioutil.TempDir("", "tempdir-*")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		return nil, nil, err
	}

	cfg, err := repo.Config()
	if err != nil || author == nil || len(author.Emails()) == 0 {
		return nil, nil, err
	}
	cfg.User.Name = author.Name()
	cfg.User.Email = author.Emails()[0]

	w, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}

	for i := 1; i <= numberOfCommits; i++ {
		_, err := w.Commit(fmt.Sprintf("Empty commit %d", i), &git.CommitOptions{
			Author: &object.Signature{
				Name:  author.Name(),
				Email: author.Emails()[0],
				When:  time.Now().Add(time.Duration(i) * time.Second),
			},
			AllowEmptyCommits: true,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return &tempDir, &cleanup, nil
}
