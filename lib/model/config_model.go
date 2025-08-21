package model

var CurrentVersion = "v1"

type ConfigModel struct {
	Version       string `mapstructure:"version"`
	TrackedAuthor struct {
		Name   string   `mapstructure:"name"`
		Emails []string `mapstructure:"emails"`
	} `mapstructure:"author"`
	TrackedRepos []string `mapstructure:"repos"`
}


func NewEmptyConfig() *ConfigModel {
	return &ConfigModel{
		Version:       CurrentVersion,
		TrackedAuthor: struct {
			Name   string   `mapstructure:"name"`
			Emails []string `mapstructure:"emails"`
		}{},
		TrackedRepos: []string{},
	}
}
