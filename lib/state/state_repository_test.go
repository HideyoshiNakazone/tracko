package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/repo"
)

func Test_NewTrackedRepo(t *testing.T) {
	repoPath := "/path/to/repo"
	lastCommit := "abc123"

	repo := NewTrackedRepo(repoPath, lastCommit)

	if repo.RepoPath != repoPath {
		t.Errorf("Expected RepoPath to be %s, got %s", repoPath, repo.RepoPath)
	}
	if repo.LastScanned != nil {
		t.Error("Expected LastScanned to be set, got zero value")
	}
}

func Test_NewCommitState(t *testing.T) {
	test_repo := NewTrackedRepo("/path/to/repo", "abc123")

	name := "John Doe"
	email := "john.doe@example.com"
	CommitID := "abc123"
	message := "Initial commit"
	commitDate := time.Now()

	state := test_repo.NewCommitState(
		name,
		email,
		CommitID,
		message,
		commitDate,
	)

	if state.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, state.Name)
	}
	if state.Email != email {
		t.Errorf("Expected Email to be %s, got %s", email, state.Email)
	}
	if state.CommitID != CommitID {
		t.Errorf("Expected CommitID to be %s, got %s", CommitID, state.CommitID)
	}
	if state.Message != message {
		t.Errorf("Expected Message to be %s, got %s", message, state.Message)
	}
	if state.CommitDate != commitDate {
		t.Errorf("Expected CommitDate to be %v, got %v", commitDate, state.CommitDate)
	}
	if state.Exported != false {
		t.Errorf("Expected Exported to be false, got %v", state.Exported)
	}
}

func Test_CommitState_MarkExported(t *testing.T) {
	test_repo := NewTrackedRepo("/path/to/repo", "abc123")
	state := test_repo.NewCommitState("John Doe", "john.doe@example.com", "abc123", "Initial commit", time.Now())
	state.MarkExported()
	if !state.Exported {
		t.Errorf("Expected Exported to be true, got %v", state.Exported)
	}
}

func Test_StateRepository_List(t *testing.T) {
	numberOfCommits := 10

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	for i := 0; i < numberOfCommits; i++ {
		commit := trackedRepo.NewCommitState(
			"John Doe",
			"john.doe@example.com",
			fmt.Sprintf("commit%d", i),
			"Initial commit",
			time.Now(),
		)
		if err := repo.AddCommit(commit); err != nil {
			t.Fatalf("Failed to create commit: %v", err)
		}
	}

	commits, err := repo.ListCommits()
	if err != nil {
		t.Fatalf("Failed to list commits: %v", err)
	}

	if len(commits) != numberOfCommits {
		t.Errorf("Expected %d commits, got %d", numberOfCommits, len(commits))
	}
}

func Test_StateRepository_GetTrackedRepoByPath_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	fetched, err := repo.GetTrackedRepoByPath("/non/existent/path")
	if err == nil {
		t.Fatal("Expected error when getting non-existent tracked repo, got nil")
	}

	if fetched != nil {
		t.Errorf("Expected nil tracked repo, got %v", fetched)
	}
}

func Test_StateRepository_CreateAndGetCommitByID(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	test_repo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.db.Create(test_repo).Error; err != nil {
		t.Fatalf("Failed to create tracked repo: %v", err)
	}

	commit := test_repo.NewCommitState("John Doe", "john.doe@example.com", "abc123", "Initial commit", time.Now())
	if err := repo.AddCommit(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	fetched, err := repo.GetCommitByID(commit.Id)
	if err != nil {
		t.Fatalf("Failed to get commit by ID: %v", err)
	}

	if fetched == nil {
		t.Fatal("Expected non-nil commit")
	}
	if fetched.Id != commit.Id {
		t.Errorf("Expected ID to be %d, got %d", commit.Id, fetched.Id)
	}
}

func Test_StateRepository_GetCommitByID_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	test_repo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.db.Create(test_repo).Error; err != nil {
		t.Fatalf("Failed to create tracked repo: %v", err)
	}

	commit := test_repo.NewCommitState("John Doe", "john.doe@example.com", "abc123", "Initial commit", time.Now())
	if err := repo.AddCommit(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	fetched, err := repo.GetCommitByID(999)
	if err == nil {
		t.Fatal("Expected error when getting non-existent commit, got nil")
	}
	if fetched != nil {
		t.Errorf("Expected nil commit, got %v", fetched)
	}
}

func Test_StateRepository_BulkAddCommit(t *testing.T) {
	numberOfCommits := 5

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	var commits []*CommitState
	for i := 0; i < numberOfCommits; i++ {
		commit := trackedRepo.NewCommitState(
			"John Doe",
			"john.doe@example.com",
			fmt.Sprintf("commit%d", i),
			"Initial commit",
			time.Now(),
		)
		commits = append(commits, commit)
	}

	if err := repo.BulkAddCommits(commits); err != nil {
		t.Fatalf("Failed to bulk add commits: %v", err)
	}
}

func Test_StateRepository_ListCommits(t *testing.T) {
	numberOfCommits := 10

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	for i := 0; i < numberOfCommits; i++ {
		commit := trackedRepo.NewCommitState(
			"John Doe",
			"john.doe@example.com",
			fmt.Sprintf("commit%d", i),
			"Initial commit",
			time.Now(),
		)
		err := repo.AddCommit(commit)
		if err != nil {
			t.Fatalf("Failed to add commit: %v", err)
		}
	}

	commits, err := repo.ListCommits()
	if err != nil {
		t.Fatalf("Failed to list commits: %v", err)
	}

	if len(commits) != numberOfCommits {
		t.Errorf("Expected 2 commits, got %d", len(commits))
	}
}

// Test TableName methods
func Test_TrackedRepo_TableName(t *testing.T) {
	repo := &TrackedRepo{}
	expected := "tracked_repo"
	if repo.TableName() != expected {
		t.Errorf("Expected TableName to be %s, got %s", expected, repo.TableName())
	}
}

func Test_CommitState_TableName(t *testing.T) {
	commit := &CommitState{}
	expected := "commit_state"
	if commit.TableName() != expected {
		t.Errorf("Expected TableName to be %s, got %s", expected, commit.TableName())
	}
}

// Test NewCommitStateFromMetadata
func Test_TrackedRepo_NewCommitStateFromMetadata(t *testing.T) {
	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	trackedRepo.Id = 1 // Set an ID for testing

	commitDate := time.Now()
	metadata := &repo.GitCommitMeta{
		AuthorName:  "Jane Doe",
		AuthorEmail: "jane.doe@example.com",
		CommitID:    "def456",
		CommitDate:  commitDate,
		Message:     "Test commit",
	}

	commitState := trackedRepo.NewCommitStateFromMetadata(metadata)

	if commitState.Name != metadata.AuthorName {
		t.Errorf("Expected Name to be %s, got %s", metadata.AuthorName, commitState.Name)
	}
	if commitState.Email != metadata.AuthorEmail {
		t.Errorf("Expected Email to be %s, got %s", metadata.AuthorEmail, commitState.Email)
	}
	if commitState.CommitID != metadata.CommitID {
		t.Errorf("Expected CommitID to be %s, got %s", metadata.CommitID, commitState.CommitID)
	}
	if commitState.CommitDate != metadata.CommitDate {
		t.Errorf("Expected CommitDate to be %v, got %v", metadata.CommitDate, commitState.CommitDate)
	}
	if commitState.Message != metadata.Message {
		t.Errorf("Expected Message to be %s, got %s", metadata.Message, commitState.Message)
	}
	if commitState.Exported != false {
		t.Errorf("Expected Exported to be false, got %v", commitState.Exported)
	}
	if commitState.TrackedRepoID != trackedRepo.Id {
		t.Errorf("Expected TrackedRepoID to be %d, got %d", trackedRepo.Id, commitState.TrackedRepoID)
	}
}

// Test StateRepository.GetTrackedRepoByID
func Test_StateRepository_GetTrackedRepoByID_Success(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	fetched, err := repo.GetTrackedRepoByID(trackedRepo.Id)
	if err != nil {
		t.Fatalf("Failed to get tracked repo by ID: %v", err)
	}

	if fetched == nil {
		t.Fatal("Expected non-nil tracked repo")
	}
	if fetched.Id != trackedRepo.Id {
		t.Errorf("Expected ID to be %d, got %d", trackedRepo.Id, fetched.Id)
	}
	if fetched.RepoPath != trackedRepo.RepoPath {
		t.Errorf("Expected RepoPath to be %s, got %s", trackedRepo.RepoPath, fetched.RepoPath)
	}
}

func Test_StateRepository_GetTrackedRepoByID_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	fetched, err := repo.GetTrackedRepoByID(999)
	if err == nil {
		t.Fatal("Expected error when getting non-existent tracked repo, got nil")
	}
	if fetched != nil {
		t.Errorf("Expected nil tracked repo, got %v", fetched)
	}
}

// Test StateRepository.GetTrackedRepoByPath success case
func Test_StateRepository_GetTrackedRepoByPath_Success(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	repoPath := "/path/to/repo"
	trackedRepo := NewTrackedRepo(repoPath, "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	fetched, err := repo.GetTrackedRepoByPath(repoPath)
	if err != nil {
		t.Fatalf("Failed to get tracked repo by path: %v", err)
	}

	if fetched == nil {
		t.Fatal("Expected non-nil tracked repo")
	}
	if fetched.RepoPath != repoPath {
		t.Errorf("Expected RepoPath to be %s, got %s", repoPath, fetched.RepoPath)
	}
}

// Test StateRepository.UpdateTrackedRepo
func Test_StateRepository_UpdateTrackedRepo(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	// Update the repo path
	newPath := "/new/path/to/repo"
	trackedRepo.RepoPath = newPath

	if err := repo.UpdateTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to update tracked repo: %v", err)
	}

	// Verify the update
	fetched, err := repo.GetTrackedRepoByID(trackedRepo.Id)
	if err != nil {
		t.Fatalf("Failed to get tracked repo after update: %v", err)
	}

	if fetched.RepoPath != newPath {
		t.Errorf("Expected updated RepoPath to be %s, got %s", newPath, fetched.RepoPath)
	}
}

// Test StateRepository.ListTrackedRepos
func Test_StateRepository_ListTrackedRepos(t *testing.T) {
	numberOfRepos := 3

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	// Add multiple tracked repos
	for i := 0; i < numberOfRepos; i++ {
		trackedRepo := NewTrackedRepo(fmt.Sprintf("/path/to/repo%d", i), fmt.Sprintf("commit%d", i))
		if err := repo.AddTrackedRepo(trackedRepo); err != nil {
			t.Fatalf("Failed to add tracked repo %d: %v", i, err)
		}
	}

	repos, err := repo.ListTrackedRepos()
	if err != nil {
		t.Fatalf("Failed to list tracked repos: %v", err)
	}

	if len(repos) != numberOfRepos {
		t.Errorf("Expected %d tracked repos, got %d", numberOfRepos, len(repos))
	}
}

// Test StateRepository.GetLastRepoCommit
func Test_StateRepository_GetLastRepoCommit_Success(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	// Add commits with different dates
	earlierTime := time.Now().Add(-2 * time.Hour)
	laterTime := time.Now().Add(-1 * time.Hour)

	commit1 := trackedRepo.NewCommitState("John Doe", "john@example.com", "commit1", "First commit", earlierTime)
	commit2 := trackedRepo.NewCommitState("Jane Doe", "jane@example.com", "commit2", "Second commit", laterTime)

	if err := repo.AddCommit(commit1); err != nil {
		t.Fatalf("Failed to add first commit: %v", err)
	}
	if err := repo.AddCommit(commit2); err != nil {
		t.Fatalf("Failed to add second commit: %v", err)
	}

	lastCommit, err := repo.GetLastRepoCommit(trackedRepo.Id)
	if err != nil {
		t.Fatalf("Failed to get last repo commit: %v", err)
	}

	if lastCommit == nil {
		t.Fatal("Expected non-nil last commit")
	}
	if lastCommit.CommitID != "commit2" {
		t.Errorf("Expected last commit ID to be 'commit2', got %s", lastCommit.CommitID)
	}
	if !lastCommit.CommitDate.Equal(laterTime) {
		t.Errorf("Expected last commit date to be %v, got %v", laterTime, lastCommit.CommitDate)
	}
}

func Test_StateRepository_GetLastRepoCommit_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	// Try to get last commit for non-existent repo
	lastCommit, err := repo.GetLastRepoCommit(999)
	if err == nil {
		t.Fatal("Expected error when getting last commit for non-existent repo, got nil")
	}
	if lastCommit != nil {
		t.Errorf("Expected nil last commit, got %v", lastCommit)
	}
}

// Test StateRepository.GetCommitCount
func Test_StateRepository_GetCommitCount(t *testing.T) {
	numberOfCommits := 5

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	// Initially should have 0 commits
	count, err := repo.GetCommitCount()
	if err != nil {
		t.Fatalf("Failed to get initial commit count: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected initial commit count to be 0, got %d", count)
	}

	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo: %v", err)
	}

	// Add commits
	for i := 0; i < numberOfCommits; i++ {
		commit := trackedRepo.NewCommitState(
			"John Doe",
			"john.doe@example.com",
			fmt.Sprintf("commit%d", i),
			"Test commit",
			time.Now(),
		)
		if err := repo.AddCommit(commit); err != nil {
			t.Fatalf("Failed to add commit %d: %v", i, err)
		}
	}

	count, err = repo.GetCommitCount()
	if err != nil {
		t.Fatalf("Failed to get commit count: %v", err)
	}
	if count != int64(numberOfCommits) {
		t.Errorf("Expected commit count to be %d, got %d", numberOfCommits, count)
	}
}

// Test StateRepository.AddTrackedRepo with duplicates
func Test_StateRepository_AddTrackedRepo_Duplicate(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	repoPath := "/path/to/repo"
	trackedRepo1 := NewTrackedRepo(repoPath, "abc123")
	trackedRepo2 := NewTrackedRepo(repoPath, "def456") // Same path, different lastCommit

	// Add first repo
	if err := repo.AddTrackedRepo(trackedRepo1); err != nil {
		t.Fatalf("Failed to add first tracked repo: %v", err)
	}

	// Add second repo with same path (should not create duplicate due to OnConflict)
	if err := repo.AddTrackedRepo(trackedRepo2); err != nil {
		t.Fatalf("Failed to add second tracked repo: %v", err)
	}

	// Verify only one repo exists
	repos, err := repo.ListTrackedRepos()
	if err != nil {
		t.Fatalf("Failed to list tracked repos: %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("Expected 1 tracked repo after duplicate insert, got %d", len(repos))
	}
}

// Test NewStateRepository with invalid path
func Test_NewStateRepository_InvalidPath(t *testing.T) {
	// Try to create repository with invalid path (empty string should work for SQLite)
	// But let's test with an invalid SQLite URL
	_, err := NewStateRepository("invalid://path")
	if err == nil {
		t.Error("Expected error when creating repository with invalid path, got nil")
	}
}

// Test that migration creates both tables correctly
func Test_NewStateRepository_Migration(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	// Test that both tables were created by performing operations that require both
	trackedRepo := NewTrackedRepo("/path/to/repo", "abc123")
	if err := repo.AddTrackedRepo(trackedRepo); err != nil {
		t.Fatalf("Failed to add tracked repo (migration issue): %v", err)
	}

	commit := trackedRepo.NewCommitState("John Doe", "john@example.com", "abc123", "Test", time.Now())
	if err := repo.AddCommit(commit); err != nil {
		t.Fatalf("Failed to add commit (migration issue): %v", err)
	}

	// If we get here, both tables exist and work correctly
}
