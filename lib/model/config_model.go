package model

var CurrentVersion = "v1"
var DefaultDBPath = "$HOME/.config/tracko.db"

type ConfigAuthorModel struct {
	Name   string   `mapstructure:"name"`
	Emails []string `mapstructure:"emails"`
}

type ConfigModel struct {
	Version       string            `mapstructure:"version"`
	DBPath        string            `mapstructure:"db_path"`
	TrackedAuthor ConfigAuthorModel `mapstructure:"author"`
	TrackedRepos  []string          `mapstructure:"repos"`
}

type ConfigModelBuilder struct {
	config *ConfigModel
}

func NewConfigBuilder() *ConfigModelBuilder {
	return &ConfigModelBuilder{
		config: &ConfigModel{
			Version:       CurrentVersion,
			DBPath:        DefaultDBPath,
			TrackedAuthor: ConfigAuthorModel{},
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

func (c *ConfigModelBuilder) WithTrackedRepos(repos []string) *ConfigModelBuilder {
	c.config.TrackedRepos = repos
	return c
}

func (c *ConfigModelBuilder) Build() *ConfigModel {
	return c.config
}
