package config

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// EnvKeyReplacer replace for environment variable parse
var EnvKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

func setup() {
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
}

// GetConfigEnvironment read config from environment variables and config.toml file
func GetConfigEnvironment(ctx context.Context, appName string) (*Config, error) {
	once.Do(func() {

		setup()

		configError = viper.ReadInConfig()
		if configError != nil {
			slog.Error("Error to read configs: ", slog.String("error", configError.Error()))
			return
		}

		// set default variables on viper
		// for example, if we have a file .toml like this:
		// [db]
		// pass
		// then the ley will be DB_PASS and if we find this key on environment variables, we will set the value on viper
		for _, k := range viper.AllKeys() {
			key := strings.ToUpper(EnvKeyReplacer.Replace(k))
			envValue := os.Getenv(key)
			if envValue != "" {
				viper.Set(k, envValue) // set as default (ignoring config file value)
			}
		}

		config = &Config{
			ctx:     ctx,
			appName: appName,
		}
		configError = viper.Unmarshal(config)
		if configError != nil {
			slog.Error("Error to unmarshal configs: ", slog.String("error", configError.Error()))
			return
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			if in.Op == fsnotify.Write {
				err := viper.Unmarshal(config)
				if err != nil {
					slog.Error("Error to unmarshal new config changes: ", slog.String("error", err.Error()))
					return
				}
			}
		})
	})

	return config, configError
}
