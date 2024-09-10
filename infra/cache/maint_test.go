package cache

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
)

var (
	testRedis *CacheManager
	cfg       *configmock.ConfigMock = configmock.New()
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	close := SetRedisTestContainerConfig(ctx, cfg)
	defer close()

	redis, client, err := NewRedisCache(ctx,
		cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.DefaultExpiration,
		cfg.GetLogger(),
	)
	if err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
	}
	defer client.Close()

	testRedis = redis

	os.Exit(m.Run())
}
