package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/diegoclair/apperr"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type dbConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type scanner interface {
	Scan(dest ...any) error
}

// handleDBError converts PostgreSQL-specific errors to apperr errors.
//   - pgx.ErrNoRows → apperr.ErrRecordNotFound
//   - unique_violation (23505) → apperr.ErrDuplicateEntry
//   - others → returned as-is
func handleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrRecordNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperr.ErrDuplicateEntry
	}

	return err
}

// queries is the connection a repository reads through. Repositories embed it,
// so the helpers below read as methods and each call carries one argument less.
type queries struct {
	db dbConn
}

// queryOne runs q and scans the single row it returns. handleDBError is applied
// here, so a caller must not map the result again.
func (c queries) queryOne[T any](ctx context.Context, q string, scan func(scanner) (T, error), args ...any) (T, error) {
	item, err := scan(c.db.QueryRow(ctx, q, args...))
	return item, handleDBError(err)
}

// queryList runs q and collects one T per row. On success the slice is never
// nil, so a caller can encode an empty result as [] instead of null. Every
// failure — the query, a scan, the iteration itself — goes through
// handleDBError. A paginated read passes withCount(q) and lets its scan close
// over the total.
func (c queries) queryList[T any](ctx context.Context, q string, scan func(scanner) (T, error), args ...any) ([]T, error) {
	rows, err := c.db.Query(ctx, q, args...)
	if err != nil {
		return nil, handleDBError(err)
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, handleDBError(err)
		}
		items = append(items, item)
	}

	return items, handleDBError(rows.Err())
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
// Example: buildInPlaceholders([]any{1}, []string{"a", "b"}, 2) returns ("$2, $3", []any{1, "a", "b"})
func buildInPlaceholders[T any](args []any, values []T, startIndex int) (placeholders string, newArgs []any) {
	if len(values) == 0 {
		return "", args
	}

	parts := make([]string, len(values))
	newArgs = args
	for i := range values {
		parts[i] = fmt.Sprintf("$%d", startIndex+i)
		newArgs = append(newArgs, values[i])
	}

	placeholders = strings.Join(parts, ", ")
	return placeholders, newArgs
}
