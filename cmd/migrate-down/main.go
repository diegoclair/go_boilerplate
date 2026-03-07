package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"

	pgMigrator "github.com/diegoclair/go_boilerplate/migrator/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		envOr("DB_POSTGRES_USERNAME", "root"),
		envOr("DB_POSTGRES_PASSWORD", "root"),
		envOr("DB_POSTGRES_HOST", "localhost"),
		envOr("DB_POSTGRES_PORT", "5432"),
		envOr("DB_POSTGRES_DB_NAME", "go_boilerplate_db"),
		envOr("DB_POSTGRES_SSLMODE", "disable"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	sqlFS, err := fs.Sub(pgMigrator.SqlFiles, "sql")
	if err != nil {
		log.Fatalf("failed to get sql subdirectory: %v", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, sqlFS)
	if err != nil {
		log.Fatalf("failed to create goose provider: %v", err)
	}

	result, err := provider.Down(ctx)
	if err != nil {
		log.Fatalf("migration down failed: %v", err)
	}

	if result == nil {
		fmt.Println("no migrations to roll back")
		return
	}

	fmt.Printf("rolled back: %s (took %s)\n", result.Source.Path, result.Duration)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
