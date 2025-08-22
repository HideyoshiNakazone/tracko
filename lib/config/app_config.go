package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/HideyoshiNakazone/tracko/lib/internal_errors"
	"github.com/HideyoshiNakazone/tracko/lib/model"
)

var trackedPaths = []string{
	"$HOME/.config/tracko", // First priority
	".",                    // Second priority
}

const configFormat string = "yaml"


func PrepareConfig(filePath string) error {
	if filePath == "" {
		for _, path := range trackedPaths {
			viper.AddConfigPath(path)
		}
	} else {
		viper.SetConfigFile(filePath)
	}

	viper.SetConfigType(configFormat)
	// other configurations be placed here, like migrations
	_, err := GetConfig()

	return err
}


func GetConfig() (*model.ConfigModel, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, internal_errors.ErrConfigNotInitialized
	}
	
	var cfg model.ConfigModel
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Version == "" {
		return nil, internal_errors.ErrConfigNotInitialized
	}

	return &cfg, nil
}

func SetConfig(cfg *model.ConfigModel) error {
	m := map[string]any{}
	if err := mapstructure.Decode(cfg, &m); err != nil {
		return err
	}
	if err := viper.MergeConfigMap(m); err != nil {
		return err
	}
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
