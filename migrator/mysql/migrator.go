package mysql

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var SqlFiles embed.FS

// Migrate applies all pending database migrations using goose.
func Migrate(db *sql.DB) error {
	ctx := context.Background()

	sqlFS, err := fs.Sub(SqlFiles, "sql")
	if err != nil {
		return fmt.Errorf("getting sql subdirectory: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectMySQL, db, sqlFS)
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
