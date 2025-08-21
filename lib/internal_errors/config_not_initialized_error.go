package internal_errors

import (
	"fmt"

	"github.com/spf13/viper"
)


type ConfigNotInitializedError struct {
	configPath string
}


func NewConfigNotInitializedError() *ConfigNotInitializedError {
	return &ConfigNotInitializedError{
		configPath: viper.ConfigFileUsed(),
	}
}

func (e *ConfigNotInitializedError) Error() string {
    return fmt.Sprintf("error: %s: config not initialized", e.configPath)
}
