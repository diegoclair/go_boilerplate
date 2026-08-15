package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn    *PostgresConn
	onceDB  sync.Once
	connErr error
)

// PostgresConn is the database connection manager
type PostgresConn struct {
	pool *pgxpool.Pool

	accountRepo contract.AccountRepo
	authRepo    contract.AuthRepo
}

// Instance returns an instance of a PostgresConn
func Instance(ctx context.Context, dsn string, log logger.Logger) (*PostgresConn, *pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	onceDB.Do(func() {
		log.Info(ctx, "Connecting to database...")

		pool, connErr = pgxpool.New(ctx, dsn)
		if connErr != nil {
			log.Error(ctx, "Database connection error", logger.Err(connErr))
			return
		}

		log.Info(ctx, "Database Ping...")
		connErr = pool.Ping(ctx)
		if connErr != nil {
			log.Error(ctx, "Database Ping error", logger.Err(connErr))
			return
		}

		log.Info(ctx, "Database successfully configured")

		conn = repoInstances(pool)
		conn.pool = pool
	})

	if conn != nil {
		pool = conn.pool
	}

	return conn, pool, connErr
}

func repoInstances(db dbConn) *PostgresConn {
	return &PostgresConn{
		accountRepo: newAccountRepo(db),
		authRepo:    newAuthRepo(db),
	}
}

// The callback gets contract.Repos, never the DataManager: a transaction opened
// from inside this one would run on its own connection, commit on its own, and
// survive the outer rollback. Handing over the repositories alone makes that a
// compile error instead of a convention.
func (c *PostgresConn) WithTransaction(ctx context.Context, fn func(tx contract.Repos) error) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	// Covers the panic path, where neither branch below runs and the connection
	// would stay checked out with the transaction open. No-op after a commit.
	defer func() { _ = tx.Rollback(ctx) }()

	// repoInstances carries no pool on purpose: nothing reachable from here may
	// start a transaction.
	txConn := repoInstances(tx)
	err = fn(txConn)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (c *PostgresConn) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *PostgresConn) Account() contract.AccountRepo {
	return c.accountRepo
}

func (c *PostgresConn) Auth() contract.AuthRepo {
	return c.authRepo
}
