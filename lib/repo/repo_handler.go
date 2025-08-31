package repo

import (
	"github.com/go-git/go-git/v6"
)


func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}
