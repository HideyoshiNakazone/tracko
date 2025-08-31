package config_model

import "github.com/HideyoshiNakazone/tracko/lib/internal_errors"

// ConfigModelBuilder is a builder for ConfigModel
type ConfigModelBuilder struct {
	config *ConfigModel
}

func NewConfigBuilder() *ConfigModelBuilder {
	return &ConfigModelBuilder{
		config: &ConfigModel{
			Version:       CurrentVersion,
			DBPath:        DefaultDBPath,
			TrackedAuthor: ConfigAuthorModel{},
			TargetRepo:    "",
			TrackedRepos:  []string{},
		},
	}
}

func (c *ConfigModelBuilder) WithDBPath(dbPath string) *ConfigModelBuilder {
	c.config.DBPath = dbPath
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthor(name string, emails []string) *ConfigModelBuilder {
	c.config.TrackedAuthor.Name = name
	c.config.TrackedAuthor.Emails = emails
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthorName(name string) *ConfigModelBuilder {
	c.config.TrackedAuthor.Name = name
	return c
}

func (c *ConfigModelBuilder) WithTrackedAuthorEmails(emails []string) *ConfigModelBuilder {
	c.config.TrackedAuthor.Emails = emails
	return c
}

func (c *ConfigModelBuilder) WithAppendTrackedAuthorEmail(email string) *ConfigModelBuilder {
	c.config.TrackedAuthor.Emails = append(c.config.TrackedAuthor.Emails, email)
	return c
}

func (c *ConfigModelBuilder) WithTargetRepo(repo string) *ConfigModelBuilder {
	c.config.TargetRepo = repo
	return c
}

func (c *ConfigModelBuilder) WithTrackedRepos(repos []string) *ConfigModelBuilder {
	c.config.TrackedRepos = repos
	return c
}

func (c *ConfigModelBuilder) Build() (*ConfigModel, error) {
	if c.config.Version == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.DBPath == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.TrackedAuthor.Name == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	if len(c.config.TrackedAuthor.Emails) == 0 {
		return nil, internal_errors.ErrInvalidConfig
	}

	if c.config.TargetRepo == "" {
		return nil, internal_errors.ErrInvalidConfig
	}

	return c.config, nil
}
