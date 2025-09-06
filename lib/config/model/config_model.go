package config_model

import (
	"fmt"
	"slices"
)

var CurrentVersion = "v1"
var DefaultDBPath = "$HOME/.config/tracko.db"

// Internal Config Model
// These struct should be completely immutable
type ConfigAuthorModel struct {
	name   string   `mapstructure:"name"`
	emails []string `mapstructure:"emails"`
}

func (a ConfigAuthorModel) Name() string {
	return a.name
}

func (a ConfigAuthorModel) Emails() []string {
	return a.emails
}

type ConfigModel struct {
	version       string
	dbPath        string
	trackedAuthor ConfigAuthorModel
	targetRepo    string
	trackedRepos  []string
}

// Getters for config
func (c ConfigModel) Version() string {
	return c.version
}

func (c ConfigModel) DBPath() string {
	return c.dbPath
}

func (c ConfigModel) TrackedAuthor() ConfigAuthorModel {
	return c.trackedAuthor
}

func (c ConfigModel) TargetRepo() string {
	return c.targetRepo
}

func (c ConfigModel) TrackedRepos() []string {
	return c.trackedRepos
}

// Manipulation methods for config
func (c ConfigModel) AppendTrackedRepo(repo string) (*ConfigModel, error) {
	repoIndex := slices.Index(c.trackedRepos, repo)
	if repoIndex != -1 {
		return nil, fmt.Errorf("repo %s already exists", repo)
	}
	c.trackedRepos = append(c.trackedRepos, repo)
	return &c, nil
}

func (c ConfigModel) RemoveTrackedRepo(repo string) (*ConfigModel, error) {
	repoIndex := slices.Index(c.trackedRepos, repo)
	if repoIndex == -1 {
		return nil, fmt.Errorf("repo %s not found", repo)
	}
	c.trackedRepos = slices.Delete(c.trackedRepos, repoIndex, repoIndex+1)
	return &c, nil
}

// External Config DTO
type AuthorDTO struct {
	Name   string   `mapstructure:"name"`
	Emails []string `mapstructure:"emails"`
}

func (a AuthorDTO) ToModel() *ConfigAuthorModel {
	return &ConfigAuthorModel{
		name:   a.Name,
		emails: a.Emails,
	}
}

type ConfigDTO struct {
	Version       string    `mapstructure:"version" restricted:"true"`
	DBPath        string    `mapstructure:"db_path"`
	TrackedAuthor AuthorDTO `mapstructure:"author"`
	TargetRepo    string    `mapstructure:"target_repo"`
	TrackedRepos  []string  `mapstructure:"tracked_repos"`
}

func (c ConfigDTO) ToModel() (*ConfigModel, error) {
	trackedAuthor := c.TrackedAuthor.ToModel()
	if trackedAuthor == nil {
		return nil, fmt.Errorf("invalid author")
	}
	return &ConfigModel{
		version:       c.Version,
		dbPath:        c.DBPath,
		trackedAuthor: *trackedAuthor,
		targetRepo:    c.TargetRepo,
		trackedRepos:  c.TrackedRepos,
	}, nil
}

func ConfigDTOFromModel(model *ConfigModel) (*ConfigDTO, error) {
	if model == nil {
		return nil, fmt.Errorf("invalid config model")
	}
	return &ConfigDTO{
		Version: model.version,
		DBPath:  model.dbPath,
		TrackedAuthor: AuthorDTO{
			Name:   model.trackedAuthor.name,
			Emails: model.trackedAuthor.emails,
		},
		TargetRepo:   model.targetRepo,
		TrackedRepos: model.trackedRepos,
	}, nil
}
