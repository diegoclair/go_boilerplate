package postgres

import (
	"testing"
)

func TestMigrate_NilPool(t *testing.T) {
	// Goose requires a real PostgreSQL connection for full migration testing.
	// The actual integration test lives in infra/data/postgres/main_test.go
	// using testcontainers with a real PostgreSQL instance.
	// This test ensures the function returns an error with a nil pool.
	err := Migrate(nil)
	if err == nil {
		t.Fatal("expected error with nil pool, got nil")
	}
}
