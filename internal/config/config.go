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
	setDefault()
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

func setDefault() {
	Config.SetDefault("host", "0.0.0.0")
	Config.SetDefault("port", 3000)
	Config.SetDefault("log.level", "info")
	Config.SetDefault("log.file", "main.log")
	//How often does the cleaning process run
	Config.SetDefault("clear_cron", 120)
	//How much is the site allowed to be monitored without polling
	Config.SetDefault("clear_interval", 3600)
	//How often are sites polled
	Config.SetDefault("interval", 55)
	//Timeout for every monitoring request
	Config.SetDefault("timeout", 30)
	//User agent for sites requests
	Config.SetDefault("user_agent", "Mozilla/5.0 (compatible; WebMonitoring/1.0;)")
}
