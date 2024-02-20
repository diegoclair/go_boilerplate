package cache

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils-lib/v2/logger"
)

// Instance returns a CacheManager instance
func Instance(ctx context.Context, cfg *config.Config, log logger.Logger) (contract.CacheManager, error) {
	return NewRedisCache(ctx, cfg, log)
}
