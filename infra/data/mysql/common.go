package mysql

import (
	"context"
	"database/sql"
)

type dbConn interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type scanner interface {
	Scan(dest ...any) error
}

func getTotalRecordsFromQuery(ctx context.Context, db dbConn, query string, args ...any) (totalRecords int64, err error) {
	var queryCount = `
		SELECT COUNT(*) FROM (
	` + query + `) as count`

	stmt, err := db.PrepareContext(ctx, queryCount)
	if err != nil {
		return totalRecords, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)

	err = row.Scan(&totalRecords)
	if err != nil {
		return totalRecords, err
	}

	return totalRecords, nil
}
