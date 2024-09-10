package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/data/mysql"
	"github.com/diegoclair/go_utils/logger"
)

// Connect returns a instance of mysql db
func Connect(ctx context.Context,
	host, port, username, password, dbName string,
	log logger.Logger,
) (contract.DataManager, *sql.DB, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting root dir: %w", err)
	}

	migrationsDir := filepath.Join(rootDir, "infra/data/migrations/mysql")

	return mysql.Instance(ctx, host, port, username, password, dbName, log, migrationsDir)
}
