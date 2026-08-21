package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/diegoclair/apperr"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type fakeRow struct{ err error }

func (r fakeRow) Scan(...any) error { return r.err }

type fakeRows struct {
	remaining int
	scanErr   error
	iterErr   error
}

func (r *fakeRows) Next() bool {
	if r.remaining == 0 {
		return false
	}
	r.remaining--
	return true
}
func (r *fakeRows) Scan(...any) error                            { return r.scanErr }
func (r *fakeRows) Err() error                                   { return r.iterErr }
func (r *fakeRows) Close()                                       {}
func (r *fakeRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *fakeRows) Values() ([]any, error)                       { return nil, nil }
func (r *fakeRows) RawValues() [][]byte                          { return nil }
func (r *fakeRows) Conn() *pgx.Conn                              { return nil }

type fakeConn struct {
	row      pgx.Row
	rows     pgx.Rows
	queryErr error
	gotQuery string
}

func (c *fakeConn) QueryRow(_ context.Context, sql string, _ ...any) pgx.Row {
	c.gotQuery = sql
	return c.row
}

func (c *fakeConn) Query(_ context.Context, sql string, _ ...any) (pgx.Rows, error) {
	c.gotQuery = sql
	return c.rows, c.queryErr
}

func (c *fakeConn) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

// Calls Scan so a driver failure reaches the helper, the way a real scan does.
func scanInt(row scanner) (int, error) {
	if err := row.Scan(new(int)); err != nil {
		return 0, err
	}
	return 7, nil
}

func TestQueryOne(t *testing.T) {
	ctx := context.Background()

	t.Run("Should return the scanned value", func(t *testing.T) {
		q := queries{db: &fakeConn{row: fakeRow{}}}

		got, err := q.queryOne(ctx, "SELECT 7", scanInt)

		require.NoError(t, err)
		require.Equal(t, 7, got)
	})

	t.Run("Should map no rows to not found", func(t *testing.T) {
		q := queries{db: &fakeConn{row: fakeRow{err: pgx.ErrNoRows}}}

		_, err := q.queryOne(ctx, "SELECT 7", scanInt)

		require.ErrorIs(t, err, apperr.ErrRecordNotFound)
	})
}

func TestQueryList(t *testing.T) {
	ctx := context.Background()

	t.Run("Should collect one value per row", func(t *testing.T) {
		q := queries{db: &fakeConn{rows: &fakeRows{remaining: 3}}}

		got, err := q.queryList(ctx, "SELECT 7", scanInt)

		require.NoError(t, err)
		require.Equal(t, []int{7, 7, 7}, got)
	})

	// The caller encodes this straight to JSON; a nil slice would ship null.
	t.Run("Should return an empty slice, never nil", func(t *testing.T) {
		q := queries{db: &fakeConn{rows: &fakeRows{}}}

		got, err := q.queryList(ctx, "SELECT 7", scanInt)

		require.NoError(t, err)
		require.NotNil(t, got)
		require.Empty(t, got)
	})

	// rows.Err is the only place a result set truncated mid-read shows up.
	// Without it a partial page is served as if it were the whole page.
	t.Run("Should map a failed iteration", func(t *testing.T) {
		q := queries{db: &fakeConn{rows: &fakeRows{remaining: 1, iterErr: errors.New("connection reset")}}}

		_, err := q.queryList(ctx, "SELECT 7", scanInt)

		require.Error(t, err)
	})

	t.Run("Should map a failed query", func(t *testing.T) {
		q := queries{db: &fakeConn{queryErr: errors.New("syntax error")}}

		_, err := q.queryList(ctx, "SELECT 7", scanInt)

		require.Error(t, err)
	})
}

// A paginated read passes withCount(q) and closes over the total, so the page
// and its count cost one round trip instead of a second COUNT query.
func TestQueryListWithCount(t *testing.T) {
	ctx := context.Background()
	db := &fakeConn{rows: &fakeRows{remaining: 2}}
	q := queries{db: db}

	var total int64
	items, err := q.queryList(ctx, withCount("SELECT id FROM tab_test"), func(row scanner) (int, error) {
		total = 42
		return scanInt(row)
	})

	require.NoError(t, err)
	require.Contains(t, db.gotQuery, "COUNT(*) OVER() as total_count")
	require.Len(t, items, 2)
	require.Equal(t, int64(42), total)
}
