package mysql

import (
	"database/sql"
)

type dbConnection interface {
	Prepare(query string) (*sql.Stmt, error)
}

type scanner interface {
	Scan(dest ...interface{}) error
}
