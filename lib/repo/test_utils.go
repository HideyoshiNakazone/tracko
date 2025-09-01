//go:build test
// +build test

package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func PrepareTestRepo(author *config_model.ConfigAuthorModel, numberOfCommits int) (*string, *func(), error) {
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
