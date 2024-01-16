package config

import (
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// Profile configuration profile
type Profile string

func (p Profile) String() string {
	return string(p)
}

// Profile default values
const (
	ProfileRun  Profile = "config"
	ProfileTest Profile = "config" // if you need a different file for test, you can create a different and change the name here
)

var (
	config      *Config
	configError error
	once        sync.Once
)

// EnvKeyReplacer replace for environment variable parse
var EnvKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

func setup(profile Profile) {
	viper.AutomaticEnv()

	viper.SetConfigName(profile.String())
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if profile == ProfileTest {
		// a test file that use config, can be in deep structure, so we need to add more paths to be possible find the config file.
		viper.AddConfigPath("../")
		viper.AddConfigPath("../../")
		viper.AddConfigPath("../../../")
		viper.AddConfigPath("../../../../")
	}
}

// GetConfigEnvironment to read initial config from a config file with it full path
func GetConfigEnvironment(profile Profile) (*Config, error) {
	once.Do(func() {

		setup(profile)

		configError = viper.ReadInConfig()
		if configError != nil {
			log.Error("Error to read configs: ", configError)
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
