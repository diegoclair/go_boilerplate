package config

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/cache"
	infraContract "github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/infra/data/mysql"
	infraLogger "github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/util/crypto"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/redis/go-redis/v9"
)

var (
	authToken infraContract.AuthToken
	authOnce  sync.Once
)

// GetAuthToken returns a new auth token or panics if it fails
func (c *Config) GetAuthToken() infraContract.AuthToken {
	authOnce.Do(func() {
		var (
			err error
			log logger.Logger = c.GetLogger()
		)

		authToken, err = auth.NewAuthToken(
			c.App.Auth.AccessTokenDuration,
			c.App.Auth.RefreshTokenDuration,
			c.App.Auth.PasetoSymmetricKey,
			log,
		)
		if err != nil {
			log.Fatalf(c.ctx, "Failed to create auth token", err)
		}
	})

	return authToken
}

var (
	cacheManager contract.CacheManager
	cacheOnce    sync.Once
)

// GetCacheManager returns a new cache manager or panics if it fails
func (c *Config) GetCacheManager() contract.CacheManager {
	cacheOnce.Do(func() {
		var (
			client *redis.Client
			log    logger.Logger = c.GetLogger()
			err    error
		)

		log.Infof(c.ctx, "Connecting to the cache server at %s:%d.", c.Cache.Redis.Host, c.Cache.Redis.Port)
		cacheManager, client, err = cache.NewRedisCache(c.ctx,
			fmt.Sprintf("%s:%d", c.Cache.Redis.Host, c.Cache.Redis.Port),
			c.Cache.Redis.Pass,
			c.Cache.Redis.DB,
			c.Cache.Redis.DefaultExpiration,
			log)
		if err != nil {
			log.Fatalf(c.ctx, "Failed to create cache manager", err)
		}

		c.AddCloser(func() {
			log.Info(c.ctx, "Closing redis connection...")
			if err := client.Close(); err != nil {
				log.Errorf(c.ctx, "Error closing redis connection: %v", err)
			}
		})
	})

	return cacheManager
}

var (
	cryptoClient contract.Crypto
	cryptoOnce   sync.Once
)

// GetCrypto returns a new crypto or panics if it fails
func (c *Config) GetCrypto() contract.Crypto {
	cryptoOnce.Do(func() {
		cryptoClient = crypto.NewCrypto()
	})

	return cryptoClient
}

var (
	dataManager contract.DataManager
	dataOnce    sync.Once
)

// GetDataManager returns a new data manager or panics if it fails
func (c *Config) GetDataManager() contract.DataManager {
	dataOnce.Do(func() {
		var (
			err     error
			mysqlDB *sql.DB
			log     logger.Logger = c.GetLogger()
		)

		dataManager, mysqlDB, err = mysql.Instance(c.ctx,
			c.GetMysqlDsn(),
			c.DB.MySQL.DBName,
			log,
		)
		if err != nil {
			log.Fatalf(c.ctx, "Failed to create data manager", err)
		}

		c.AddCloser(func() {
			log.Info(c.ctx, "Closing mysql connection...")
			if err := mysqlDB.Close(); err != nil {
				log.Errorf(c.ctx, "Error closing mysql connection: %v", err)
			}
		})
	})

	return dataManager
}

var (
	l       logger.Logger
	logOnce sync.Once
)

// GetLogger returns a new logger
func (c *Config) GetLogger() logger.Logger {
	logOnce.Do(func() {
		l = infraLogger.NewLogger(c.appName, c.Log.Debug)
	})

	return l
}

var (
	v             validator.Validator
	validatorOnce sync.Once
)

// GetValidator returns a new validator or panics if it fails
func (c *Config) GetValidator() validator.Validator {
	validatorOnce.Do(func() {
		var (
			err error
			log logger.Logger = c.GetLogger()
		)

		v, err = validator.NewValidator()
		if err != nil {
			log.Fatalf(c.ctx, "Failed to create validator", err)
		}
	})

	return v
}

func (c *Config) GetHttpPort() string {
	return c.App.Port
}
