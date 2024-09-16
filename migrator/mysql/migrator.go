package mysql

import (
	"database/sql"
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/GuiaBolso/darwin"
)

//go:embed sql/*.sql
var SqlFiles embed.FS

// Migrate applies database migrations to the provided SQL database.
// It reads SQL migration files from the specified directory and executes them in order.
// The function returns an error if any migration fails to execute.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})
	path := "sql"

	files, err := SqlFiles.ReadDir(path)
	if err != nil {
		return err
	}

	migrations := []darwin.Migration{}

	for _, file := range files {
		migrationFile := path + "/" + file.Name()

		migration, err := SqlFiles.ReadFile(migrationFile)
		if err != nil {
			return err
		}

		migrations = append(migrations, getMigrations(file.Name(), migration)...)

	}

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}

func getMigrations(filename string, migration []byte) []darwin.Migration {
	migrationString := string(migration)
	instructions := strings.Split(migrationString, ";")

	result := []darwin.Migration{}
	fileVersion := strings.Split(filename, "_")[0]

	// the darwin package can't handle multiple instructions in the same file, so we need to split the file into multiple instructions(migrations)
	for i, instruction := range instructions {
		version, err := getFileVersion(fileVersion, i)
		if err != nil {
			panic(err)
		}

		// the last instruction is empty because the last character is ";" and the split function will return an empty string
		if strings.TrimSpace(instruction) == "" {
			continue
		}

		action := strings.TrimSpace(strings.Split(instruction, "(")[0])

		result = append(result, darwin.Migration{
			Version:     version,
			Description: fmt.Sprintf("%s; %s", action, filename),
			Script:      instruction + ";",
		})
	}

	return result
}

func getFileVersion(fileVersion string, i int) (float64, error) {
	fv, err := strconv.Atoi(fileVersion)
	if err != nil {
		return 0, err
	}

	stringVersion := fmt.Sprintf("%d.%05d", fv, i)
	version, err := strconv.ParseFloat(stringVersion, 64)
	if err != nil {
		return 0, err
	}

	return version, nil
}
