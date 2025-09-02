package state

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewCommitState(t *testing.T) {
	name := "John Doe"
	email := "john.doe@example.com"
	repoPath := "/path/to/repo"
	CommitID := "abc123"
	message := "Initial commit"
	commitDate := time.Now()

	state := NewCommitState(name, email, repoPath, CommitID, message, commitDate)

	if state.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, state.Name)
	}
	if state.Email != email {
		t.Errorf("Expected Email to be %s, got %s", email, state.Email)
	}
	if state.RepoPath != repoPath {
		t.Errorf("Expected RepoPath to be %s, got %s", repoPath, state.RepoPath)
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
	state := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	state.MarkExported()
	if !state.Exported {
		t.Errorf("Expected Exported to be true, got %v", state.Exported)
	}
}

func Test_NewStateRepository(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}
	if repo == nil {
		t.Fatal("Expected non-nil state repository")
	}
}

func Test_StateRepository_CreateAndGetByID(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	commit := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	if err := repo.Create(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	fetched, err := repo.GetByID(commit.Id)
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

func Test_StateRepository_GetByID_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	commit := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	if err := repo.Create(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	fetched, err := repo.GetByID(999)
	if err == nil {
		t.Fatal("Expected error when getting non-existent commit, got nil")
	}
	if fetched != nil {
		t.Errorf("Expected nil commit, got %v", fetched)
	}
}

func Test_StateRepository_List(t *testing.T) {
	numberOfCommits := 10

	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	for i := 0; i < numberOfCommits; i++ {
		commit := NewCommitState(
			"John Doe",
			"john.doe@example.com",
			"/path/to/repo",
			fmt.Sprintf("commit%d", i),
			"Initial commit",
			time.Now(),
		)
		if err := repo.Create(commit); err != nil {
			t.Fatalf("Failed to create commit: %v", err)
		}
	}

	commits, err := repo.List()
	if err != nil {
		t.Fatalf("Failed to list commits: %v", err)
	}

	if len(commits) != numberOfCommits {
		t.Errorf("Expected 2 commits, got %d", len(commits))
	}
}

func Test_StateRepository_Update(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	commit := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	if err := repo.Create(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	commit.Message = "Updated commit message"
	if err := repo.Update(commit); err != nil {
		t.Fatalf("Failed to update commit: %v", err)
	}

	fetched, err := repo.GetByID(commit.Id)
	if err != nil {
		t.Fatalf("Failed to get commit by ID: %v", err)
	}

	if fetched == nil {
		t.Fatal("Expected non-nil commit")
	}
	if fetched.Message != commit.Message {
		t.Errorf("Expected Message to be %s, got %s", commit.Message, fetched.Message)
	}
}

func Test_StateRepository_Update_NoID(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	commit := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	if err := repo.Create(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	commit.Id = 0
	if err := repo.Update(commit); err == nil {
		t.Fatal("Expected error when updating commit with no ID, got nil")
	}
}

func Test_StateRepository_Delete(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	commit := NewCommitState("John Doe", "john.doe@example.com", "/path/to/repo", "abc123", "Initial commit", time.Now())
	if err := repo.Create(commit); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	if err := repo.Delete(commit.Id); err != nil {
		t.Fatalf("Failed to delete commit: %v", err)
	}

	fetched, err := repo.GetByID(commit.Id)
	if err == nil {
		t.Fatal("Expected error when getting deleted commit, got nil")
	}
	if fetched != nil {
		t.Errorf("Expected nil commit, got %v", fetched)
	}
}

func Test_StateRepository_Delete_NotFound(t *testing.T) {
	repo, err := NewStateRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create state repository: %v", err)
	}

	if err := repo.Delete(999); err == nil {
		t.Fatal("Expected error when deleting non-existent commit, got nil")
	}
}
