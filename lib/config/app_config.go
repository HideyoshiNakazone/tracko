package config

import (
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/HideyoshiNakazone/tracko/lib/model"
)

var trackedPaths = []string{
	"$HOME/.config/tracko", // First priority
	"/etc/tracko",          // Second priority
	".",                    // Third priority
}



const configFormat string = "yaml"
var configInitialized bool = false



func prepareConfig() error {
	viper.SetConfigType(configFormat)

	if configInitialized {
		return nil
	}
	configInitialized = true

	if err := viper.ReadInConfig(); err != nil {
		SetConfig(
			model.NewEmptyConfig(),
		)
	}

	_, err := GetConfig()
	return err
}



func InitializeConfig() error {
	for _, path := range trackedPaths {
		viper.AddConfigPath(path)
	}

	err := prepareConfig()
	if err != nil {
		return err
	}

	return nil
}


func InitializeConfigFromFile(filePath string) error {
	viper.SetConfigFile(filePath)

	err := prepareConfig()
	if err != nil {
		return err
	}

	return nil
}



func GetConfig() (*model.ConfigModel, error) {
	if !configInitialized {
		return nil, errors.New("Config not initialized")
	}

	var cfg model.ConfigModel
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}


func SetConfig(cfg *model.ConfigModel) error {
	if !configInitialized {
		return errors.New("Config not initialized")
	}

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
