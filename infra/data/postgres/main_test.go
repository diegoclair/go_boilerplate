package postgres

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	pgMigrator "github.com/diegoclair/go_boilerplate/migrator/postgres"
)

var (
	testDB contract.DataManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg := configmock.New()

	close := setPostgresTestContainerConfig(ctx, cfg)
	defer close()

	pgConn, pool, err := Instance(ctx,
		cfg.GetPostgresDSN(),
		cfg.GetLogger(),
	)
	if err != nil {
		log.Fatalf("cannot connect to postgres: %v", err)
	}

	err = pgMigrator.Migrate(pool)
	if err != nil {
		log.Fatalf("cannot migrate postgres: %v", err)
	}

	testDB = pgConn

	os.Exit(m.Run())
}
