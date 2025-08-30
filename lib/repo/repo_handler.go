package repo

import (
	"path/filepath"

	"github.com/go-git/go-git/v6"
)


func IsGitRepository(path *string) (string, bool) {
	absPath, err := filepath.Abs(*path)
	if err != nil {
		return "", false
	}

	_, err = git.PlainOpen(absPath)
	if err != nil {
		return "", false
	}
	return absPath, true
}
