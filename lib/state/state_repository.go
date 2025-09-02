package state

import (
	"errors"
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CommitState struct {
	Id         uint      `gorm:"primaryKey"`             // Primary key
	Name       string    `gorm:"size:100;not null"`      // VARCHAR(100), NOT NULL
	Email      string    `gorm:"size:100;not null"`      // VARCHAR(100), NOT NULL
	RepoPath   string    `gorm:"size:255;not null"`      // VARCHAR(255), NOT NULL
	CommitID   string    `gorm:"size:40;not null"`       // VARCHAR(40), NOT NULL
	CommitDate time.Time `gorm:"not null"`               // TIMESTAMP, NOT NULL
	Message    string    `gorm:"type:text;not null"`     // TEXT, NOT NULL
	Exported   bool      `gorm:"default:false;not null"` // BOOLEAN, NOT NULL
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
	return r.db.Create(commit).Error
}

// GetByID finds a CommitState by ID.
func (r *StateRepository) GetByID(id uint) (*CommitState, error) {
	var commit CommitState
	if err := r.db.First(&commit, id).Error; err != nil {
		return nil, err
	}
	return &commit, nil
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
