package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig[T any](configName string) (T, error) {
	var config T

	confPath, err := filepath.Abs("./configs")
	if err != nil {
		return config, err
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
