package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var SqlFiles embed.FS

// Migrate applies all pending database migrations using goose.
func Migrate(pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("pool is nil")
	}

	ctx := context.Background()

	db := stdlib.OpenDBFromPool(pool)

	sqlFS, err := fs.Sub(SqlFiles, "sql")
	if err != nil {
		return fmt.Errorf("getting sql subdirectory: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, sqlFS)
	if err != nil {
		return fmt.Errorf("creating goose provider: %w", err)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	for _, r := range results {
		log.Printf("[goose] applied migration: %s (%s, %v)", r.Source.Path, r.Direction, r.Duration)
	}

	return nil
}
