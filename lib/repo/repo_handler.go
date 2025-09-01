package repo

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// IsGitRepository checks if the given path is a Git repository.
func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

// Commit metadata
type GitCommitMeta struct {
	AuthorName  string
	AuthorEmail string
	RepoPath    string
	CommitID    string
	CommitDate  time.Time
	Message     string
}

func GitCommitMetaFromObject(
	repoPath string,
	commit *object.Commit,
) (*GitCommitMeta, error) {
	if commit == nil {
		return nil, errors.New("commit is nil")
	}

	return &GitCommitMeta{
		AuthorName:  commit.Author.Name,
		AuthorEmail: commit.Author.Email,
		RepoPath:    repoPath,
		CommitID:    commit.ID().String(),
		CommitDate:  commit.Author.When,
		Message:     commit.Message,
	}, nil
}

// Commit Iterator
type CommitIter interface {
	Next() (*GitCommitMeta, error)
	ForEach(func(*GitCommitMeta) error) error
	Close()
}

// commitIterator implements object.CommitIter.
type commitIterator struct {
	repoPath string
	iter     object.CommitIter
	filters  func(*GitCommitMeta) bool
}

// Next advances the iterator and returns the next matching commit.
func (it *commitIterator) Next() (*GitCommitMeta, error) {
	commit, err := it.iter.Next()
	if err != nil {
		return nil, err
	}

	meta, err := GitCommitMetaFromObject(it.repoPath, commit)
	if err != nil {
		return nil, err
	}

	if !it.filters(meta) {
		return it.Next()
	}
	return meta, nil
}

func (it *commitIterator) ForEach(fn func(*GitCommitMeta) error) error {
	for {
		commit, err := it.Next()
		if err != nil {
			return err
		}
		if commit == nil {
			break
		}
		if err := fn(commit); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the underlying commit iterator.
func (it *commitIterator) Close() {
	it.iter.Close()
}

// Checks that commitIterator implements CommitIter
var _ CommitIter = &commitIterator{}

// TrackedRepo represents a Git repository being tracked.
type TrackedRepo struct {
	repoPath string
	repo     *git.Repository
	author   *config_model.ConfigAuthorModel
}

func NewTrackedRepo(path string, author *config_model.ConfigAuthorModel) (*TrackedRepo, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &TrackedRepo{
		repoPath: path,
		repo:     gitRepo,
		author:   author,
	}, nil
}

type ListRepositoryHistoryParams struct {
	Since *time.Time
	Until *time.Time
}

func (r *TrackedRepo) ListRepositoryHistory(options *ListRepositoryHistoryParams) (CommitIter, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repository not initialized")
	}

	if options == nil {
		options = &ListRepositoryHistoryParams{}
	}

	iter, err := r.repo.Log(&git.LogOptions{
		Since: options.Since,
		Until: options.Until,
	})
	if err != nil {
		return nil, err
	}

	filter := func(meta *GitCommitMeta) bool {
		if r.author == nil || meta.AuthorName == "" {
			// Filters by Author by default, therefore these values are needed
			return false
		}

		if r.author.Name() != meta.AuthorName {
			return false
		}
		if !slices.Contains(r.author.Emails(), meta.AuthorEmail) {
			return false
		}

		return true
	}

	return &commitIterator{
		iter:     iter,
		filters:  filter,
		repoPath: r.repoPath,
	}, nil
}
