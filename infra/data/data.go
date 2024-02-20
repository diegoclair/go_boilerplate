package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data/mysql"
	"github.com/diegoclair/go_utils-lib/v2/logger"
)

// Connect returns a instance of mysql db
func Connect(ctx context.Context, cfg *config.Config, log logger.Logger) (contract.DataManager, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting root dir: %w", err)
	}

	migrationsDir := filepath.Join(rootDir, "infra/data/migrations/mysql")

	return mysql.Instance(ctx, cfg, log, migrationsDir)
}
