package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	Config *viper.Viper
)

func LoadOrCreatePersistentConfig(configFilePath string) error {
	Config = viper.New()
	Config.SetConfigType("json")
	Config.SetConfigFile(configFilePath)

	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := Config.SafeWriteConfig(); err != nil {
				return errors.New(fmt.Sprintf("Couldn't create config file: %v", err))
			}
		} else {
			return errors.New(fmt.Sprintf("Couldn't read config file: %v", err))
		}
	}

	return nil
}
