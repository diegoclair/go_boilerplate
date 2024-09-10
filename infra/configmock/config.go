package configmock

import (
	"fmt"
	"time"

	"github.com/diegoclair/go_utils/logger"
)

type ConfigMock struct {
	Redis RedisConfig
	Auth  AuthConfig
	DB    DBConfig
}

type AuthConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PasetoSymmetricKey   string
}

type DBConfig struct {
	MySQL MySQLConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host              string
	Password          string
	DB                int
	DefaultExpiration time.Duration
}

func New() *ConfigMock {
	return &ConfigMock{
		Auth: AuthConfig{
			AccessTokenDuration:  time.Minute * 15,
			RefreshTokenDuration: time.Hour * 24,
			PasetoSymmetricKey:   "d152a3402-4d32-85ad-19df4c9934cd",
		},
		DB: DBConfig{
			MySQL: MySQLConfig{
				Host:     "",
				Port:     "",
				Username: "guest",
				Password: "guest",
				DBName:   "test",
			},
		},
		Redis: RedisConfig{
			Host:              "",
			Password:          "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
			DB:                0,
			DefaultExpiration: time.Minute,
		},
	}
}

func (c *ConfigMock) GetLogger() logger.Logger {
	return logger.NewNoop()
}

func (c *ConfigMock) SetRedisHost(host, port string) {
	c.Redis.Host = fmt.Sprintf("%s:%s", host, port)
}

func (c *ConfigMock) SetMySQLHostAndPort(host, port string) {
	c.DB.MySQL.Host = host
	c.DB.MySQL.Port = port
}
