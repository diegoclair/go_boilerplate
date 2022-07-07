package cache

import (
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/util/config"
)

// Instance returns a CacheManager instance
func Instance(cfg config.RedisConfig, log logger.Logger) (contract.CacheManager, error) {
	return NewRedisCache(cfg, log)
}
