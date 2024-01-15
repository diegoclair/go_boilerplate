package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/GuiaBolso/darwin"
)

// Migrate applies database migrations to the provided SQL database.
// It reads SQL migration files from the specified directory and executes them in order.
// The function returns an error if any migration fails to execute.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})

	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting root dir: %w", err)
	}

	migrationsDir := filepath.Join(rootDir, "infra/data/migrations/mysql")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	migrations := []darwin.Migration{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFile := migrationsDir + "/" + file.Name()

			migration, err := os.ReadFile(migrationFile)
			if err != nil {
				return err
			}

			migrations = append(migrations, getMigrations(file.Name(), migration)...)
		}
	}

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}

func getMigrations(filename string, migration []byte) []darwin.Migration {
	migrationString := string(migration)
	instructions := strings.Split(migrationString, ";")

	result := []darwin.Migration{}

	// the darwin package can't handle multiple instructions in the same file, so we need to split the file into multiple instructions(migrations)
	for i, instruction := range instructions {
		version, err := getFileVersion(filename, i)
		if err != nil {
			panic(err)
		}

		// the last instruction is empty because the last character is ";" and the split function will return an empty string
		if strings.TrimSpace(instruction) == "" {
			continue
		}

		result = append(result, darwin.Migration{
			Version:     version,
			Description: filename,
			Script:      instruction + ";",
		})
	}

	return result

}

func getFileVersion(filename string, i int) (float64, error) {
	version := strings.Split(filename, "_")[0]
	v, err := strconv.Atoi(version)
	if err != nil {
		return 0, err
	}

	instructionPosition := fmt.Sprintf("0.%d", i)
	position, err := strconv.ParseFloat(instructionPosition, 64)
	if err != nil {
		return 0, err
	}

	return float64(v) + position, nil
}
