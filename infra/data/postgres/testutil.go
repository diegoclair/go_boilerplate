package postgres

import (
	"context"
	"log"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// setPostgresTestContainerConfig set the postgres container for testing
//
// You can use this function to set the postgres container for an integration testing
func setPostgresTestContainerConfig(ctx context.Context, cfg *configmock.ConfigMock) (closeFunc func()) {
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18",
		postgres.WithDatabase(cfg.DB.Postgres.DBName),
		postgres.WithUsername(cfg.DB.Postgres.Username),
		postgres.WithPassword(cfg.DB.Postgres.Password),
	)
	if err != nil {
		log.Fatalf("cannot start postgres container: %v", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	cfg.SetPostgresHostAndPort(host, port.Port())

	return func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}
}
