package state

import (
	"time"

	"github.com/HideyoshiNakazone/tracko/lib/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// TrackedRepo represents a tracked repository in the database.
type TrackedRepo struct {
	Id          uint          `gorm:"primaryKey"`               // Primary key
	RepoPath    string        `gorm:"size:255;not null;unique"` // VARCHAR(255), NOT NULL, unique
	LastScanned *time.Time     `gorm:"autoUpdateTime"`  			// TIMESTAMP, NULLABLE, auto-updated on each save
	Commits     []CommitState `gorm:"foreignKey:TrackedRepoID"` // One-to-many relation
}

func (TrackedRepo) TableName() string {
    return "tracked_repo"
}

func NewTrackedRepo(repoPath, lastCommit string) *TrackedRepo {
	return &TrackedRepo{
		RepoPath:    repoPath,
		LastScanned: nil,
	}
}

func (r *TrackedRepo) UpdateLastScanned(t time.Time) {
	r.LastScanned = &t
}

func (r *TrackedRepo) NewCommitState(name, email, commitId, message string, commitDate time.Time) *CommitState {
	return &CommitState{
		Name:          name,
		Email:         email,
		CommitID:      commitId,
		CommitDate:    commitDate,
		Message:       message,
		Exported:      false,
		TrackedRepoID: r.Id,
		TrackedRepo:   *r,
	}
}

func (r *TrackedRepo) NewCommitStateFromMetadata(metadata *repo.GitCommitMeta) *CommitState {
	return &CommitState{
		Name:       metadata.AuthorName,
		Email:      metadata.AuthorEmail,
		CommitID:   metadata.CommitID,
		CommitDate: metadata.CommitDate,
		Message:    metadata.Message,
		Exported:   false,
		TrackedRepoID: r.Id,
		TrackedRepo:   *r,
	}
}


// CommitState represents the state of a commit in the database.
type CommitState struct {
	Id            uint        `gorm:"primaryKey"`                               // Primary key
	Name          string      `gorm:"size:100;not null"`                        // VARCHAR(100), NOT NULL
	Email         string      `gorm:"size:100;not null"`                        // VARCHAR(100), NOT NULL
	CommitID      string      `gorm:"size:40;not null;uniqueIndex:idx_repo_commit_id"` // VARCHAR(40), NOT NULL
	CommitDate    time.Time   `gorm:"not null"`                                // TIMESTAMP, NOT NULL
	Message       string      `gorm:"type:text;not null"`                      // TEXT, NOT NULL
	Exported      bool        `gorm:"default:false;not null"`                  // BOOLEAN, NOT NULL

	TrackedRepoID uint        `gorm:"not null;uniqueIndex:idx_repo_commit_id"` // Foreign key
	TrackedRepo   TrackedRepo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` 
}

func (CommitState) TableName() string {
    return "commit_state"
}

func (c *CommitState) MarkExported() {
	c.Exported = true
}


// StateRepository provides methods to interact with the database.
type StateRepository struct {
	db *gorm.DB
}

// NewStateRepository initializes SQLite and auto-migrates the schema.
func NewStateRepository(dbPath string) (*StateRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&TrackedRepo{}, &CommitState{})
	if err != nil {
		return nil, err
	}

	return &StateRepository{db: db}, nil
}

// Inserts a new TrackedRepo.
func (r *StateRepository) AddTrackedRepo(repo *TrackedRepo) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(repo).Error
}

// GetByID finds a TrackedRepo by ID.
func (r *StateRepository) GetTrackedRepoByID(id uint) (*TrackedRepo, error) {
	var repo TrackedRepo
	if err := r.db.First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// GetByPath finds a TrackedRepo by its repository path.
func (r *StateRepository) GetTrackedRepoByPath(path string) (*TrackedRepo, error) {
	var repo TrackedRepo
	if err := r.db.Where("repo_path = ?", path).First(&repo).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// Updates an existing TrackedRepo.
func (r *StateRepository) UpdateTrackedRepo(repo *TrackedRepo) error {
	return r.db.Save(repo).Error
}

// List returns all tracked repositories.
func (r *StateRepository) ListTrackedRepos() ([]TrackedRepo, error) {
	var repos []TrackedRepo
	if err := r.db.Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}


// Inserts a new CommitState.
func (r *StateRepository) AddCommit(commit *CommitState) error {
	err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(commit).Error
	if err != nil {
		return err
	}

	commit.TrackedRepo.UpdateLastScanned(time.Now())
	if err := r.UpdateTrackedRepo(&commit.TrackedRepo); err != nil {
		return err
	}

	return nil
}

// Bulk inserts multiple CommitStates.
func (r *StateRepository) BulkAddCommits(commits []*CommitState) error {
	err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(commits).Error
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		return nil
	}

	repo := commits[0].TrackedRepo
	repo.UpdateLastScanned(time.Now())
	if err := r.UpdateTrackedRepo(&repo); err != nil {
		return err
	}

	return nil
}

// GetCommitByID finds a CommitState by ID.
func (r *StateRepository) GetCommitByID(id uint) (*CommitState, error) {
	var commit CommitState
	if err := r.db.First(&commit, id).Error; err != nil {
		return nil, err
	}
	return &commit, nil
}

// GetLastRepoCommit retrieves the last commit for a specific repository.
func (r *StateRepository) GetLastRepoCommit(id uint) (*CommitState, error) {
	var commit CommitState
	if err := r.db.Where("tracked_repo_id = ?", id).Order("commit_date desc").First(&commit).Error; err != nil {
		return nil, err
	}

	return &commit, nil
}

// Count number of saved commits
func (r *StateRepository) GetCommitCount() (int64, error) {
	var count int64
	if err := r.db.Model(&CommitState{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List returns all commits.
func (r *StateRepository) ListCommits() ([]CommitState, error) {
	var commits []CommitState
	if err := r.db.Find(&commits).Error; err != nil {
		return nil, err
	}
	return commits, nil
}
