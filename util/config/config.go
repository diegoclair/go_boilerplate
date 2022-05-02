package config

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// GetConfigEnvironment to read initial config
func GetConfigEnvironment() (*Config, error) {
	once.Do(func() {

		viper.SetConfigFile("config.toml")
		viper.AutomaticEnv()

		configError = viper.ReadInConfig()
		if configError != nil {
			log.Error("Error to read configs: ", configError)
			return
		}

		config = &Config{}
		configError = viper.Unmarshal(config)
		if configError != nil {
			log.Error("Error to unmarshal configs: ", configError)
			return
		}
	})

	return config, configError
}

type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"db"`
}

type AppConfig struct {
	Auth AuthConfig `mapstructure:"auth"`
}
type AuthConfig struct {
	JWTPrivateKey string `mapstructure:"jwt-private-key"`
}

type DBConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	DBName             string `mapstructure:"db-name"`
	CryptoKey          string `mapstructure:"crypto-key"`
	MaxLifeInMinutes   int    `mapstructure:"max-life-in-minutes"`
	MaxIdleConnections int    `mapstructure:"max-idle-connections"`
	MaxOpenConnections int    `mapstructure:"max-open-connections"`
}
