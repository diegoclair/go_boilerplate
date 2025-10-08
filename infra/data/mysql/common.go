package mysql

import (
	"context"
	"database/sql"
	"regexp"
)

type dbConn interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type scanner interface {
	Scan(dest ...any) error
}

// withCount adds COUNT(*) OVER() to a base query for pagination
// It searches for the first "FROM" keyword (case-insensitive, word boundary) and inserts the count column before it
func withCount(baseQuery string) string {
	// \b = word boundary, ensures FROM is a complete word (not from_id, perform, etc.)
	// (?i) = case-insensitive
	re := regexp.MustCompile(`(?i)\bFROM\b`)

	// Find first match
	loc := re.FindStringIndex(baseQuery)
	if loc == nil {
		return baseQuery
	}

	// Insert COUNT(*) OVER() before FROM
	beforeFrom := baseQuery[:loc[0]]
	fromAndAfter := baseQuery[loc[0]:]

	return beforeFrom + ",\n\t\tCOUNT(*) OVER() as total_count\n\t\t" + fromAndAfter
}
