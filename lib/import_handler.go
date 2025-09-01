package lib

import (
	"fmt"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
)

func processTrackedRepos(commitIter *repo.CommitIter, ch chan *repo.GitCommitMeta) {
	if commitIter == nil || ch == nil {
		return
	}
	defer func() {
		(*commitIter).Close()
	}()

	(*commitIter).ForEach(func(meta *repo.GitCommitMeta) error {
		ch <- meta
		return nil
	})
}

func ImportTrackedRepos(cfg *config_model.ConfigModel) error {
	author := cfg.TrackedAuthor()

	repos := make([]*repo.TrackedRepo, 0)
	for _, repoPath := range cfg.TrackedRepos() {
		repo, err := repo.NewTrackedRepo(repoPath, &author)
		if err != nil {
			return err
		}
		repos = append(repos, repo)
	}

	commitChannel := make(chan *repo.GitCommitMeta)
	for _, trackedRepo := range repos {
		commitIter, err := trackedRepo.ListRepositoryHistory(&repo.ListRepositoryHistoryParams{})
		if err != nil {
			return err
		}

		go processTrackedRepos(&commitIter, commitChannel)
	}

	go func() {
		for commit := range commitChannel {
			fmt.Println(commit)
		}
		close(commitChannel)
	}()

	return nil
}
