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
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	DBName             string `mapstructure:"db-name"`
	SSLMode            string `mapstructure:"sslmode"`
	MaxLifeInMinutes   int    `mapstructure:"max-life-in-minutes"`
	MaxIdleConnections int    `mapstructure:"max-idle-connections"`
	MaxOpenConnections int    `mapstructure:"max-open-connections"`
}

// GetPostgresDsn builds the connection string with pgxpool sizing params.
// Without pool_max_conns pgx defaults to max(4, NumCPU), ignoring the
// configured limit.
func (c *Config) GetPostgresDsn() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB.Postgres.Username, c.DB.Postgres.Password,
		c.DB.Postgres.Host, c.DB.Postgres.Port,
		c.DB.Postgres.DBName, c.DB.Postgres.SSLMode,
	)
	if c.DB.Postgres.MaxOpenConnections > 0 {
		dsn += fmt.Sprintf("&pool_max_conns=%d", c.DB.Postgres.MaxOpenConnections)
	}
	if c.DB.Postgres.MaxIdleConnections > 0 {
		dsn += fmt.Sprintf("&pool_min_conns=%d", c.DB.Postgres.MaxIdleConnections)
	}
	if c.DB.Postgres.MaxLifeInMinutes > 0 {
		dsn += fmt.Sprintf("&pool_max_conn_lifetime=%dm", c.DB.Postgres.MaxLifeInMinutes)
	}
	return dsn
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
