package mysql

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// setMysqlTestContainerConfig set the mysql container for testing
//
// You can use this function to set the mysql container for an integration testing
func setMysqlTestContainerConfig(ctx context.Context, cfg *config.Config) (closeFunc func()) {
	mysqlContainer, err := mysql.Run(
		ctx,
		"mysql:8.0.32",
		mysql.WithDatabase(cfg.DB.MySQL.DBName),
		mysql.WithUsername(cfg.DB.MySQL.Username),
		mysql.WithPassword(cfg.DB.MySQL.Password),
	)
	if err != nil {
		log.Fatalf("cannot start mysql container: %v", err)
	}

	cfg.DB.MySQL.Host, err = mysqlContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	cfg.DB.MySQL.Port = port.Port()

	return func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}
}

// testForUpdateDeleteErrorsWithMock is a helper function to test the update and delete functions
func testForUpdateDeleteErrorsWithMock(t *testing.T, f func(db *sql.DB) error) {
	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			err = f(db)
			require.Error(t, err)
		})
	}
}

// testForInsertErrorsWithMock is a helper function to test the insert functions
func testForInsertErrorsWithMock(t *testing.T, f func(db *sql.DB) error) {
	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the last insert id fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().
					WillReturnResult(sqlmock.NewErrorResult(assert.AnError))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			err = f(db)
			require.Error(t, err)
		})
	}
}

// testForSelectErrorsWithMock is a helper function to test the select functions
func testForSelectErrorsWithMock(t *testing.T, tableID string, f func(db *sql.DB) error) {
	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{tableID}).
						AddRow("invalid"), // id of table should be an integer
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			err = f(db)
			require.Error(t, err)
		})
	}
}

// testForPaginatedSelectErrorsWithMock is a helper function to test the paginated select functions
func testForPaginatedSelectErrorsWithMock(t *testing.T, tableID string, f func(db *sql.DB) error) {
	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "Should not return error when the count returns 0",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(0),
				)
			},
		},
		{
			name: "Should return error when the count prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the count query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the address scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{tableID}).
						AddRow("invalid"), // id of table should be an integer
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			err = f(db)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
