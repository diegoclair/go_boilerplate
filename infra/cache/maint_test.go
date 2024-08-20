package cache

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
)

var (
	testRedis contract.CacheManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg, err := config.GetConfigEnvironment(config.ProfileTest)
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	close := SetRedisTestContainerConfig(ctx, cfg)
	defer close()

	redis, err := Instance(ctx, cfg, logger.NewNoop())
	if err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
	}

	testRedis = redis

	os.Exit(m.Run())
}
