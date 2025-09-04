package import_handler

import (
	"fmt"
	"sync"

	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"github.com/HideyoshiNakazone/tracko/lib/state"
	"github.com/HideyoshiNakazone/tracko/lib/utils"
)

func processTrackedRepos(repoPath string, state_repo *state.StateRepository, cfg *config_model.ConfigModel, ch chan *repo.GitCommitMeta, wg *sync.WaitGroup) {
	defer wg.Done()

	author := cfg.TrackedAuthor()
	trackedRepo, err := repo.NewTrackedRepo(repoPath, &author)
	if err != nil {
		return
	}

	listParams := repo.ListRepositoryHistoryParams{}
	lastCommit, err := state_repo.GetLastRepoCommit(repoPath)
	if err == nil {
		listParams.Since = &lastCommit.CommitDate
	} else {
		fmt.Println("No last commit found, importing all commits")
		err = nil
	}

	commitIter, err := trackedRepo.ListRepositoryHistory(&listParams)
	if err != nil {
		return
	}
	defer commitIter.Close()

	commitIter.ForEach(func(meta *repo.GitCommitMeta) error {
		ch <- meta
		return nil
	})
}

func processCommits(state_repo *state.StateRepository, ch chan *repo.GitCommitMeta) {
	batchSize := 1_000

	for batch := range utils.PartitionChannel(ch, batchSize) {
		err := state_repo.BulkCreate(
			utils.Map(batch, state.NewCommitStateFromMetadata),
		)
		if err != nil {
			fmt.Printf("Failed to bulk create commits: %v\n", err)
			return
		}
	}

	commitCount, err := state_repo.Count()
	if err != nil {
		fmt.Printf("Failed to count commits: %v\n", err)
		return
	}
	fmt.Printf("Imported %d commits\n", commitCount)
}

func ImportTrackedRepos(cfg *config_model.ConfigModel) error {
	commitChannel := make(chan *repo.GitCommitMeta)
	state_repo, err := state.NewStateRepository(cfg.DBPath())
	if err != nil {
		return err
	}

	var read_wg sync.WaitGroup
	for _, repoPath := range cfg.TrackedRepos() {
		read_wg.Add(1)
		go processTrackedRepos(repoPath, state_repo, cfg, commitChannel, &read_wg)
	}

	// Start the writer goroutine (no WaitGroup needed)
	done := make(chan struct{})
	go func() {
		processCommits(state_repo, commitChannel)
		close(done)
	}()

	read_wg.Wait()
	close(commitChannel)
	<-done // Wait for writer to finish

	return nil
}
