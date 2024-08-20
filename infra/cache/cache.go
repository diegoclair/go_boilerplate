package cache

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
)

// Instance returns a CacheManager instance
func Instance(ctx context.Context, cfg *config.Config, log logger.Logger) (contract.CacheManager, error) {
	return newRedisCache(ctx, cfg, log)
}
