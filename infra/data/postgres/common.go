package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/diegoclair/go_boilerplate/internal/domain"
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

// handleDBError converts PostgreSQL-specific errors to domain errors.
//   - pgx.ErrNoRows → domain.ErrNotFound
//   - unique_violation (23505) → domain.ErrConflict
//   - others → returned as-is
func handleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrConflict
	}

	return err
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
