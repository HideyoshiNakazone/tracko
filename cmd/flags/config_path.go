package flags

import "os"


var ConfigPath string



func GetConfigPath() string {
	configPath := os.Getenv("TRACKO_CONFIG_PATH")
	if configPath != "" {
		return configPath
	}
	return ConfigPath
}
