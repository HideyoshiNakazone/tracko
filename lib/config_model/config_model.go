package config_model


var CurrentVersion = "v1"
var DefaultDBPath = "$HOME/.config/tracko.db"

type ConfigAuthorModel struct {
	Name   string   `mapstructure:"name"`
	Emails []string `mapstructure:"emails"`
}

type ConfigModel struct {
	Version       string            `mapstructure:"version" restricted:"true"`
	DBPath        string            `mapstructure:"db_path"`
	TrackedAuthor ConfigAuthorModel `mapstructure:"author"`
	TargetRepo    string            `mapstructure:"target_repo"`
	TrackedRepos  []string          `mapstructure:"tracked_repos"`
}

