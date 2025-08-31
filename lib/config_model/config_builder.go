package config_model

import "github.com/HideyoshiNakazone/tracko/lib/internal_errors"

// ConfigModelBuilder is a builder for ConfigModel
type ConfigModelBuilder struct {
	config *ConfigModel
}

func NewConfigBuilder() *ConfigModelBuilder {
	return &ConfigModelBuilder{
		config: &ConfigModel{
			version:       CurrentVersion,
			dbPath:        DefaultDBPath,
			trackedAuthor: ConfigAuthorModel{},
			targetRepo:    "",
			trackedRepos:  []string{},
		},
	}
}

func (c *ConfigModelBuilder) WithDBPath(dbPath string) *ConfigModelBuilder {
	c.config.dbPath = dbPath
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthor(name string, emails []string) *ConfigModelBuilder {
	c.config.trackedAuthor.name = name
	c.config.trackedAuthor.emails = emails
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthorName(name string) *ConfigModelBuilder {
	c.config.trackedAuthor.name = name
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthorEmails(emails []string) *ConfigModelBuilder {
	c.config.trackedAuthor.emails = emails
	return c
}

func (c *ConfigModelBuilder) WithAppendTrackedAuthorEmail(email string) *ConfigModelBuilder {
	c.config.trackedAuthor.emails = append(c.config.trackedAuthor.emails, email)
	return c
}

func (c *ConfigModelBuilder) WithTargetRepo(repo string) *ConfigModelBuilder {
	c.config.targetRepo = repo
	return c
}

func (c *ConfigModelBuilder) WithTrackedRepos(repos []string) *ConfigModelBuilder {
	c.config.trackedRepos = repos
	return c
}

func (c *ConfigModelBuilder) Build() (*ConfigModel, error) {
	if c.config.version == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.dbPath == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.trackedAuthor.name == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if len(c.config.trackedAuthor.emails) == 0 {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.targetRepo == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	return c.config, nil
}
