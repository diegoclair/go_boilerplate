package mysql

import (
	"database/sql"
	"embed"

	"github.com/GuiaBolso/darwin"
	"github.com/diegoclair/sqlmigrator"
)

//go:embed sql/*.sql
var SqlFiles embed.FS

// Migrate applies database migrations using the sqlmigrator package.
// It reads SQL migration files from the embedded sql directory and executes them in order.
// The function returns an error if any migration fails to execute.
func Migrate(db *sql.DB) error {
	// Create migrator with MySQL dialect
	migrator := sqlmigrator.New(db, darwin.MySQLDialect{})

	// Execute migrations using the embedded SqlFiles and "sql" directory
	return migrator.Migrate(SqlFiles, "sql")
}
