package import_handler

import (
	"testing"

	config_handler "github.com/HideyoshiNakazone/tracko/lib/config/handler"
	config_model "github.com/HideyoshiNakazone/tracko/lib/config/model"
	"github.com/HideyoshiNakazone/tracko/lib/repo"
)

func Test_ImportTrackedRepos(t *testing.T) {
	cfg, cleanup, err := prepareTestConfig()
	if err != nil {
		t.Fatalf("Failed to prepare test config: %v", err)
	}
	defer cleanup()

	err = ImportTrackedRepos(cfg)
	if err != nil {
		t.Fatalf("Failed to import tracked repos: %v", err)
	}
}

func prepareTestConfig() (*config_model.ConfigModel, func(), error) {
	numberOfCommits := 100

	testAuthor := config_model.AuthorDTO{
		Name: "Test User",
		Emails: []string{
			"testuser@example.com",
		},
	}.ToModel()

	repoPath1, repo_cleanup1, err := repo.PrepareTestRepo(testAuthor, numberOfCommits)
	if err != nil {
		return nil, nil, err
	}

	repoPath2, repo_cleanup2, err := repo.PrepareTestRepo(testAuthor, numberOfCommits)
	if err != nil {
		return nil, nil, err
	}

	cfg, err := config_model.NewConfigBuilder().
		WithDBPath(":memory:?cache=shared").
		WithTrackedAuthor(testAuthor.Name(), testAuthor.Emails()).
		WithTargetRepo("repo1").
		WithTrackedRepos([]string{
			*repoPath1,
			*repoPath2,
		}).
		Build()

	if err != nil {
		return nil, nil, err
	}

	_, config_cleanup, err := config_handler.PrepareTestConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		(*config_cleanup)()
		(*repo_cleanup1)()
		(*repo_cleanup2)()
	}

	return cfg, cleanup, nil
}
