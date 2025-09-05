package import_handler

import (
	"sync"

	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"github.com/HideyoshiNakazone/tracko/lib/state"
	"github.com/HideyoshiNakazone/tracko/lib/utils"
)

func processTrackedRepos(repoPath string, cfg *config_model.ConfigModel, ch chan *repo.GitCommitMeta, errorChannel chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	state_repo, err := state.NewStateRepository(cfg.DBPath())
	if err != nil {
		errorChannel <- err
		return
	}

	author := cfg.TrackedAuthor()
	trackedRepo, err := repo.NewTrackedRepo(repoPath, &author)
	if err != nil {
		errorChannel <- err
		return
	}

	listParams := repo.ListRepositoryHistoryParams{}
	lastCommit, err := state_repo.GetLastRepoCommit(repoPath)
	if err == nil {
		listParams.Since = &lastCommit.CommitDate
	} else {
		err = nil
	}

	commitIter, err := trackedRepo.ListRepositoryHistory(&listParams)
	if err != nil {
		errorChannel <- err
		return
	}
	defer commitIter.Close()

	commitIter.ForEach(func(meta *repo.GitCommitMeta) error {
		ch <- meta
		return nil
	})
}

func processCommits(cfg *config_model.ConfigModel, ch chan *repo.GitCommitMeta, errorChannel chan error) {
	batchSize := 1_000

	state_repo, err := state.NewStateRepository(cfg.DBPath())
	if err != nil {
		errorChannel <- err
		return
	}

	for batch := range utils.PartitionChannel(ch, batchSize) {
		err := state_repo.BulkCreate(
			utils.Map(batch, state.NewCommitStateFromMetadata),
		)
		if err != nil {
			errorChannel <- err
			return
		}
	}
}

func ImportTrackedRepos(cfg *config_model.ConfigModel) error {
	commitChannel := make(chan *repo.GitCommitMeta)
	errorChannel := make(chan error)

	var read_wg sync.WaitGroup
	for _, repoPath := range cfg.TrackedRepos() {
		read_wg.Add(1)
		go processTrackedRepos(repoPath, cfg, commitChannel, errorChannel, &read_wg)
	}

	// Start the writer goroutine (no WaitGroup needed)
	go func() {
		processCommits(cfg, commitChannel, errorChannel)
		close(errorChannel)
	}()
	
	go func() {
		read_wg.Wait()
		close(commitChannel)
	}()

	for err := range errorChannel {
		if err != nil {
			return err
		}
	}

    return nil

}
