package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/diegoclair/go_boilerplate/migrator/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		envOr("DB_MYSQL_USERNAME", "root"),
		envOr("DB_MYSQL_PASSWORD", "root"),
		envOr("DB_MYSQL_HOST", "localhost"),
		envOr("DB_MYSQL_PORT", "3306"),
		envOr("DB_MYSQL_DB_NAME", "go_boilerplate_db"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	sqlFS, err := fs.Sub(mysql.SqlFiles, "sql")
	if err != nil {
		log.Fatalf("failed to get sql subdirectory: %v", err)
	}

	provider, err := goose.NewProvider(goose.DialectMySQL, db, sqlFS)
	if err != nil {
		log.Fatalf("failed to create goose provider: %v", err)
	}

	result, err := provider.Down(context.Background())
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
