package mysql

import (
	"context"
	"database/sql"
)

type dbConnection interface {
	Prepare(query string) (*sql.Stmt, error)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func getTotalRecordsFromQuery(ctx context.Context, db dbConnection, query string, args ...interface{}) (totalRecords int64, err error) {
	var queryCount = `
		SELECT COUNT(*) FROM (
	` + query + `) as count`

	stmt, err := db.Prepare(queryCount)
	if err != nil {
		return totalRecords, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(args...)

	err = row.Scan(&totalRecords)
	if err != nil {
		return totalRecords, err
	}

	return totalRecords, nil
}
