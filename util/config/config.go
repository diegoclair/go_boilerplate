package config

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

var (
	config *Config
	once   sync.Once
)

// GetConfigEnvironment to read initial config
func GetConfigEnvironment() *Config {
	once.Do(func() {

		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Error("Error to read configs: ", err)
			panic(err)
		}

		config = &Config{}

		config.App.Auth.PrivateKey = viper.GetString("JWT_PRIVATE_KEY")

		config.MySQL.Username = viper.GetString("DB_USER")
		config.MySQL.Password = viper.GetString("DB_PASSWORD")
		config.MySQL.Host = viper.GetString("DB_HOST")
		config.MySQL.Port = viper.GetString("DB_PORT")
		config.MySQL.DBName = viper.GetString("DB_NAME")
		config.MySQL.CryptoKey = viper.GetString("DB_CRYPTO_KEY")
		config.MySQL.MaxLifeInMinutes = viper.GetInt("MAX_LIFE_IN_MINUTES")
		config.MySQL.MaxIdleConns = viper.GetInt("MAX_IDLE_CONNS")
		config.MySQL.MaxOpenConns = viper.GetInt("MAX_OPEN_CONNS")

	})

	return config
}

type Config struct {
	App   AppConfig
	MySQL MysqlConfig
}

type AppConfig struct {
	Auth AuthConfig
}
type AuthConfig struct {
	PrivateKey string
}

type MysqlConfig struct {
	Username         string
	Password         string
	Host             string
	Port             string
	DBName           string
	CryptoKey        string
	MaxLifeInMinutes int
	MaxIdleConns     int
	MaxOpenConns     int
}
