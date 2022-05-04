package config

import (
	"os"
	"strings"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const (
	ConfigDefaultFilepath = "config.toml"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// EnvKeyReplacer replace for environment variable parse
var EnvKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

// GetConfigEnvironment to read initial config
func GetConfigEnvironment(filepath string) (*Config, error) {
	once.Do(func() {

		viper.SetConfigFile(filepath)
		viper.AutomaticEnv()

		configError = viper.ReadInConfig()
		if configError != nil {
			log.Error("Error to read configs: ", configError)
			return
		}

		for _, k := range viper.AllKeys() {
			key := strings.ToUpper(EnvKeyReplacer.Replace(k))
			envValue := os.Getenv(key)
			if envValue != "" {
				viper.Set(k, envValue)
			}
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
	Environment string     `mapstructure:"environment"`
	Port        string     `mapstructure:"port"`
	Auth        AuthConfig `mapstructure:"auth"`
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
