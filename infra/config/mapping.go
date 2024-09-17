package config

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	App      AppConfig   `mapstructure:"app"`
	Cache    CacheConfig `mapstructure:"cache"`
	DB       DBConfig    `mapstructure:"db"`
	Log      LogConfig   `mapstructure:"log"`
	closers  []func()
	closerMu sync.Mutex
	ctx      context.Context
	appName  string
}

func (c *Config) AddCloser(close func()) {
	c.closerMu.Lock()
	defer c.closerMu.Unlock()
	c.closers = append(c.closers, close)
}

func (c *Config) Close() {
	c.closerMu.Lock()
	defer c.closerMu.Unlock()

	for _, close := range c.closers {
		close()
	}
}

type AppConfig struct {
	Name        string     `mapstructure:"name"`
	Environment string     `mapstructure:"environment"`
	Port        string     `mapstructure:"port"`
	Auth        AuthConfig `mapstructure:"auth"`
}
type AuthConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"access-token-duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh-token-duration"`
	PasetoSymmetricKey   string        `mapstructure:"paseto-symmetric-key"`
}

type CacheConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
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
	MaxLifeInMinutes   int    `mapstructure:"max-life-in-minutes"`
	MaxIdleConnections int    `mapstructure:"max-idle-connections"`
	MaxOpenConnections int    `mapstructure:"max-open-connections"`
}

func (c *Config) GetMysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", c.DB.MySQL.Username, c.DB.MySQL.Password, c.DB.MySQL.Host, c.DB.MySQL.Port, c.DB.MySQL.DBName)
}

type LogConfig struct {
	Debug     bool   `mapstructure:"debug"`
	LogToFile bool   `mapstructure:"log-to-file"`
	Path      string `mapstructure:"path"`
}

type RedisConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	DB                int           `mapstructure:"db"`
	Pass              string        `mapstructure:"pass"`
	Prefix            string        `mapstructure:"prefix"`
	DefaultExpiration time.Duration `mapstructure:"default-expiration"`
}
