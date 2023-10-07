package config

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	User     string `mapstructure:"dbuser"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	Port     string `mapstructure:"port"`
}

// LoadDBConfig
func LoadDBConfig() (DBConfig, error) {
	var config DBConfig
	if err := viper.Unmarshal(&config); err != nil {
		return DBConfig{}, err
	}
	// log.Printf("[DBconfig] user: %v, pass: %v, host: %v, db: %v, port: %v", config.User, config.Password, config.Host, config.Database, config.Port)
	return config, nil
}
