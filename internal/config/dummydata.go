package config

import (
	"log"

	"github.com/spf13/viper"
)

type DummyDataConfig struct {
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Database string `mapstructure:"DB_DATABASE"`
	Port     string `mapstructure:"DB_PORT"`
}

// LoadDummyDataConfig
func LoadDummyDataConfig() (DummyDataConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_DATABASE")
	viper.BindEnv("DB_PORT")
	var config DummyDataConfig
	if err := viper.Unmarshal(&config); err != nil {
		return DummyDataConfig{}, err
	}
	log.Printf("[DBconfig] user: %v, pass: %v, host: %v, db: %v, port: %v", config.User, config.Password, config.Host, config.Database, config.Port)
	return config, nil
}
