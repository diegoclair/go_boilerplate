package mysql

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_utils/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_instance(t *testing.T) {
	ctx := context.Background()
	cfg := configmock.New()

	db, dbMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	var getTestMysql getMysql = func(string) (*sql.DB, error) {
		return db, nil
	}

	var getTestMysqlError getMysql = func(string) (*sql.DB, error) {
		return nil, assert.AnError
	}

	type args struct {
		testMysql getMysql
	}

	tests := []struct {
		name       string
		args       args
		setupTests func(dbMock sqlmock.Sqlmock, args args)
		wantErr    bool
	}{
		{
			name: "Should return error when getMysql returns error",
			args: args{
				testMysql: getTestMysqlError,
			},
			wantErr: true,
		},
		{
			name: "Should return error when db ping fails",
			args: args{
				testMysql: getTestMysql,
			},
			setupTests: func(dm sqlmock.Sqlmock, args args) {
				dm.ExpectPing().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error if create database fails",
			args: args{
				testMysql: getTestMysql,
			},
			setupTests: func(dm sqlmock.Sqlmock, args args) {
				dm.ExpectPing().WillReturnError(nil)
				dm.ExpectExec("CREATE DATABASE").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error if use database fails",
			args: args{
				testMysql: getTestMysql,
			},
			setupTests: func(dm sqlmock.Sqlmock, args args) {
				dm.ExpectPing().WillReturnError(nil)
				dm.ExpectExec("CREATE DATABASE").WillReturnResult(sqlmock.NewResult(0, 0))
				dm.ExpectExec("USE").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTests != nil {
				tt.setupTests(dbMock, tt.args)
			}

			onceDB = sync.Once{}
			mysql, db, err := instance(ctx,
				cfg.DB.MySQL.Host,
				cfg.DB.MySQL.Port,
				cfg.DB.MySQL.Username,
				cfg.DB.MySQL.Password,
				cfg.DB.MySQL.DBName,
				logger.NewNoop(),
				"",
				tt.args.testMysql,
			)
			if tt.wantErr && err != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, mysql)
			assert.NotNil(t, db)
		})
	}

}
