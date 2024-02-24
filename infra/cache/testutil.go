package cache

import (
	"context"
	"log"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetRedisTestContainerConfig set the redis container for testing
//
// You can use this function to set the redis container for an integration testing
func SetRedisTestContainerConfig(ctx context.Context, cfg *config.Config) (closeFunc func()) {

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Server initialized"),
		// this commands are the same as the command from docker-compose.yml
		Cmd: []string{"redis-server", "--save 20 1", "--requirepass", "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"},
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("cannot start mysql container: %v", err)
	}

	cfg.Cache.Redis.Host, err = redisContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	cfg.Cache.Redis.Port = port.Int()

	return func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}
}
