package state

import (
	"errors"
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommitState struct {
	Id         uint      `gorm:"primaryKey"`              							// Primary key
	Name       string    `gorm:"size:100;not null"`       							// VARCHAR(100), NOT NULL
	Email      string    `gorm:"size:100;not null"`       							// VARCHAR(100), NOT NULL
	RepoPath   string    `gorm:"size:255;not null;uniqueIndex:idx_repo_commit_id"`  // VARCHAR(255), NOT NULL
	CommitID   string    `gorm:"size:40;not null;uniqueIndex:idx_repo_commit_id"` 	// VARCHAR(40), NOT NULL
	CommitDate time.Time `gorm:"not null"`                							// TIMESTAMP, NOT NULL
	Message    string    `gorm:"type:text;not null"`      							// TEXT, NOT NULL
	Exported   bool      `gorm:"default:false;not null"`  							// BOOLEAN, NOT NULL
}

func NewCommitState(name, email, repoPath, commitId, message string, commitDate time.Time) *CommitState {
	return &CommitState{
		Name:       name,
		Email:      email,
		RepoPath:   repoPath,
		CommitID:   commitId,
		CommitDate: commitDate,
		Message:    message,
		Exported:   false,
	}
}

func NewCommitStateFromMetadata(metadata *repo.GitCommitMeta) *CommitState {
	return &CommitState{
		Name:       metadata.AuthorName,
		Email:      metadata.AuthorEmail,
		RepoPath:   metadata.RepoPath,
		CommitID:   metadata.CommitID,
		CommitDate: metadata.CommitDate,
		Message:    metadata.Message,
		Exported:   false,
	}
}

func (c *CommitState) MarkExported() {
	c.Exported = true
}

type StateRepository struct {
	db *gorm.DB
}

// NewStateRepository initializes SQLite and auto-migrates the schema.
func NewStateRepository(dbPath string) (*StateRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&CommitState{})
	if err != nil {
		return nil, err
	}

	return &StateRepository{db: db}, nil
}

// Create inserts a new CommitState.
func (r *StateRepository) Create(commit *CommitState) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(commit).Error
}

// Bulk Create inserts multiple CommitStates.
func (r *StateRepository) BulkCreate(commits []*CommitState) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(commits).Error
}

// GetByID finds a CommitState by ID.
func (r *StateRepository) GetByID(id uint) (*CommitState, error) {
	var commit CommitState
	if err := r.db.First(&commit, id).Error; err != nil {
		return nil, err
	}
	return &commit, nil
}

// GetLastRepoCommit retrieves the last commit for a specific repository.
func (r *StateRepository) GetLastRepoCommit(repoPath string) (*CommitState, error) {
	var commit CommitState
	if err := r.db.Where("repo_path = ?", repoPath).Order("commit_date desc").First(&commit).Error; err != nil {
		return nil, err
	}
	return &commit, nil
}

// Count number of saved commits
func (r *StateRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&CommitState{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List returns all commits.
func (r *StateRepository) List() ([]CommitState, error) {
	var commits []CommitState
	if err := r.db.Find(&commits).Error; err != nil {
		return nil, err
	}
	return commits, nil
}

// Update updates an existing CommitState.
func (r *StateRepository) Update(commit *CommitState) error {
	if commit.Id == 0 {
		return errors.New("commit must have an ID to update")
	}
	return r.db.Save(commit).Error
}

// Delete deletes a commit by ID.
func (r *StateRepository) Delete(id uint) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}

	return r.db.Delete(&CommitState{}, id).Error
}
