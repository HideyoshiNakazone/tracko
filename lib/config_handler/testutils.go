package config_handler

import (
	"errors"
	"os"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
)

func PrepareTestConfig(cfg *config_model.ConfigModel) (*os.File, *func(), error) {
	tempFile, err := os.CreateTemp("", "tracko_test_config_*.yaml")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	if err := PrepareConfig(tempFile.Name()); err == nil {
		return nil, nil, errors.New("config file already exists")
	}

	err = SetConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	return tempFile, &cleanup, nil
}
