package config_handler

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/HideyoshiNakazone/tracko/lib/config_model"
	"github.com/HideyoshiNakazone/tracko/lib/internal_errors"
	"github.com/HideyoshiNakazone/tracko/lib/utils"
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

func GetConfig() (*config_model.ConfigModel, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, internal_errors.ErrConfigNotInitialized
	}

	var cfg config_model.ConfigModel
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Version == "" {
		return nil, internal_errors.ErrConfigNotInitialized
	}

	return &cfg, nil
}

func GetConfigAttr[T any](key string) (T, error) {
	var zero T

	if err := viper.ReadInConfig(); err != nil {
		return zero, fmt.Errorf("failed to read config: %w", err)
	}

	val := viper.Get(key)
	if val == nil {
		return zero, fmt.Errorf("config value for %q is nil", key)
	}

	v := reflect.ValueOf(val)
	targetType := reflect.TypeOf(zero)

	if targetType == nil {
		return val.(T), nil
	}

	if !v.Type().AssignableTo(targetType) {
		return zero, fmt.Errorf("cannot cast config value for %q from %T to %v", key, val, targetType)
	}

	casted := v.Convert(targetType).Interface().(T)
	return casted, nil
}

func SetConfig(cfg *config_model.ConfigModel) error {
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

func SetConfigAttr(key string, value any) error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if !utils.CheckModelHasField(config_model.ConfigModel{}, key) {
		return fmt.Errorf("field %q does not exist", key)
	}

	if utils.CheckModelHasTag(config_model.ConfigModel{}, key, "restricted", "true") {
		return fmt.Errorf("field %q is restricted and cannot be modified", key)
	}

	viper.Set(key, value)

	var cfg config_model.ConfigModel
	if err := viper.Unmarshal(&cfg); err != nil {
		viper.ReadInConfig()
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return viper.WriteConfig()
}
