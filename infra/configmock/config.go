package configmock

import (
	"fmt"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
	Postgres PostgresConfig
}

type PostgresConfig struct {
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
			Postgres: PostgresConfig{
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

func (c *ConfigMock) GetPostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.Postgres.Username, c.DB.Postgres.Password,
		c.DB.Postgres.Host, c.DB.Postgres.Port,
		c.DB.Postgres.DBName,
	)
}

func (c *ConfigMock) GetLogger() logger.Logger {
	return logger.NewNoop()
}

func (c *ConfigMock) SetRedisHost(host, port string) {
	c.Redis.Host = fmt.Sprintf("%s:%s", host, port)
}

func (c *ConfigMock) SetPostgresHostAndPort(host, port string) {
	c.DB.Postgres.Host = host
	c.DB.Postgres.Port = port
}

func (c *ConfigMock) GetCacheManager(ctrl *gomock.Controller) *mocks.MockCacheManager {
	return mocks.NewMockCacheManager(ctrl)
}

func (c *ConfigMock) GetValidator(t *testing.T) validator.Validator {
	v, err := validator.NewValidator()
	require.NoError(t, err)

	return v
}

func (c *ConfigMock) GetCrypto(ctrl *gomock.Controller) *mocks.MockCrypto {
	return mocks.NewMockCrypto(ctrl)
}
