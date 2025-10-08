package mysql

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
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

// buildInPlaceholders creates SQL IN clause placeholders and appends values to args
// Example: buildInPlaceholders([]any{1}, []string{"a", "b"}) returns ("?, ?", []any{1, "a", "b"})
func buildInPlaceholders[T any](args []any, values []T) (placeholders string, newArgs []any) {
	if len(values) == 0 {
		return "", args
	}

	// Build placeholders and append values in a single loop
	parts := make([]string, len(values))
	newArgs = args
	for i := range values {
		parts[i] = "?"
		newArgs = append(newArgs, values[i])
	}

	placeholders = strings.Join(parts, ", ")
	return placeholders, newArgs
}
