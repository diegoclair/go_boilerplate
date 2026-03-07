package config

import (
	"context"
	"sync"

	goconfig "github.com/diegoclair/go_utils/config"
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// GetConfigEnvironment read config from environment variables and config.toml file
func GetConfigEnvironment(ctx context.Context, appName string) (*Config, error) {
	once.Do(func() {
		config, configError = goconfig.Load[Config](goconfig.Options{
			SearchPaths:  []string{".", "../", "../../"},
			WatchChanges: true,
		})
		if configError != nil {
			return
		}

		config.ctx = ctx
		config.appName = appName
		config.setupTracer()
	})

	return config, configError
}
