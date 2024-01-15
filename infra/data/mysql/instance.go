package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"sync"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

var (
	conn    *mysqlConn
	onceDB  sync.Once
	connErr error
)

// mysqlConn is the database connection manager
type mysqlConn struct {
	db *sql.DB

	accountRepo contract.AccountRepo
	authRepo    contract.AuthRepo
}

// Instance returns an instance of a MySQLRepo
func Instance(ctx context.Context, cfg *config.Config, log logger.Logger) (contract.DataManager, error) {
	onceDB.Do(func() {

		dataSourceName := fmt.Sprintf("%s:root@tcp(%s:%s)/?charset=utf8&parseTime=true",
			cfg.DB.MySQL.Username, cfg.DB.MySQL.Host, cfg.DB.MySQL.Port,
		)

		var db *sql.DB
		log.Info(ctx, "Connecting to database...")
		db, connErr = sql.Open("mysql", dataSourceName)
		if connErr != nil {
			return
		}

		log.Info(ctx, "Database Ping...")
		connErr = db.PingContext(ctx)
		if connErr != nil {
			log.Errorf(ctx, "Database Ping error: %v", connErr)
			return
		}

		log.Info(ctx, "Creating database if not exists...")
		if _, connErr = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.DB.MySQL.DBName)); connErr != nil {
			log.Errorf(ctx, "Create Database error: %v", connErr)
			return
		}

		if _, connErr = db.Exec(fmt.Sprintf("USE %s;", cfg.DB.MySQL.DBName)); connErr != nil {
			log.Errorf(ctx, "Default Database error: %v", connErr)
			return
		}

		connErr = mysqlDriver.SetLogger(log)
		if connErr != nil {
			return
		}
		log.Info(ctx, "Database successfully configured")

		log.Info(ctx, "Running the migrations")
		connErr = Migrate(db)
		if connErr != nil {
			log.Errorf(ctx, "Migrate Error: %v", connErr)
			return
		}

		log.Info(ctx, "Migrations executed")

		conn = repoInstances(db)
		conn.db = db
	})

	return conn, connErr
}

func repoInstances(dbConn dbConnection) *mysqlConn {
	return &mysqlConn{
		accountRepo: newAccountRepo(dbConn),
		authRepo:    newAuthRepo(dbConn),
	}
}

func (c *mysqlConn) WithTransaction(ctx context.Context, fn func(dm contract.DataManager) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	txConn := repoInstances(tx)
	err = fn(txConn)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

func (c *mysqlConn) Close() (err error) {
	return c.db.Close()
}

func (c *mysqlConn) Account() contract.AccountRepo {
	return c.accountRepo
}

func (c *mysqlConn) Auth() contract.AuthRepo {
	return c.authRepo
}
