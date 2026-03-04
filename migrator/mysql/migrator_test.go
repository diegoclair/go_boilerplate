package mysql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMigrate(t *testing.T) {
	// Goose requires a real MySQL connection for full migration testing.
	// The actual integration test lives in infra/data/mysql/main_test.go
	// using testcontainers with a real MySQL instance.
	// This test ensures the function does not panic with a mock DB.
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	err = Migrate(db)
	if err != nil {
		t.Logf("expected error with mock db: %v", err)
	}
}
