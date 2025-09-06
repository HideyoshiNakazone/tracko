package import_handler

import (
	"sync"

	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"github.com/HideyoshiNakazone/tracko/lib/state"
	"github.com/HideyoshiNakazone/tracko/lib/utils"
)

func processTrackedRepos(repoPath string, state_repo *state.StateRepository, cfg *config_model.ConfigModel, ch chan *state.CommitState, errorChannel chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	author := cfg.TrackedAuthor()
	repo_handler, err := repo.NewTrackedRepo(repoPath, &author)
	if err != nil {
		errorChannel <- err
		return
	}

	tracked_repo, err := state_repo.GetTrackedRepoByPath(repoPath)
	if err != nil {
		tracked_repo = state.NewTrackedRepo(repoPath, "")
		err = state_repo.AddTrackedRepo(tracked_repo)
		if err != nil {
			errorChannel <- err
			return
		}
	}

	listParams := repo.ListRepositoryHistoryParams{}
	if tracked_repo.LastScanned != nil {
		listParams.Since = tracked_repo.LastScanned
	}

	lastCommit, err := state_repo.GetLastRepoCommit(tracked_repo.Id)
	if err == nil {
		listParams.From = &lastCommit.CommitID
	}

	commitIter, err := repo_handler.ListRepositoryHistory(&listParams)
	if err != nil {
		errorChannel <- err
		return
	}
	defer commitIter.Close()

	commitIter.ForEach(func(meta *repo.GitCommitMeta) error {
		ch <- tracked_repo.NewCommitStateFromMetadata(meta)
		return nil
	})
}

func processCommits(state_repo *state.StateRepository, ch chan *state.CommitState, errorChannel chan error) {
	batchSize := 1_000

	for batch := range utils.PartitionChannel(ch, batchSize) {
		err := state_repo.BulkAddCommits(batch)
		if err != nil {
			errorChannel <- err
			return
		}
	}
}

func ImportTrackedRepos(cfg *config_model.ConfigModel) error {
	commitChannel := make(chan *state.CommitState)
	errorChannel := make(chan error)

	state_repo, err := state.NewStateRepository(cfg.DBPath())
	if err != nil {
		return err
	}

	var read_wg sync.WaitGroup
	for _, repoPath := range cfg.TrackedRepos() {
		read_wg.Add(1)
		go processTrackedRepos(repoPath, state_repo, cfg, commitChannel, errorChannel, &read_wg)
	}

	// Start the writer goroutine (no WaitGroup needed)
	go func() {
		processCommits(state_repo, commitChannel, errorChannel)
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
