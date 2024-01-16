package config

import (
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const (
	ConfigDefaultName = "config.toml"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// EnvKeyReplacer replace for environment variable parse
var EnvKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

// GetConfigEnvironment to read initial config from a config file with it full path
func GetConfigEnvironment(filepath string) (*Config, error) {
	once.Do(func() {

		viper.SetConfigFile(filepath)
		viper.AutomaticEnv()

		configError = viper.ReadInConfig()
		if configError != nil {
			log.Error("Error to read configs: ", configError)
			return
		}

		// here we will try to find the key on Environment variables and if we find, we will Set the value with viper.
		// this set of the value, will be used as default instead of it finds on your config file
		// for example, if we have a file .toml like this:
		// [db]
		// pass
		// it will mount the key as db.pass and will transform it to DB_PASS to try to find it on env vars.
		// it helps when we have a file with local config but to deploy we haver another env var defined.
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

		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			if in.Op == fsnotify.Write {
				err := viper.Unmarshal(config)
				if err != nil {
					log.Error("Error to unmarshal new config changes: ", err)
					return
				}
			}
		})
	})

	return config, configError
}
