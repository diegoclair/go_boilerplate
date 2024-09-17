package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"sync"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/resterrors"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

var (
	conn    *MysqlConn
	onceDB  sync.Once
	connErr error
)

// MysqlConn is the database connection manager
type MysqlConn struct {
	db *sql.DB

	accountRepo contract.AccountRepo
	authRepo    contract.AuthRepo
}

// helps test the Instance function
type getMysql func(string) (*sql.DB, error)

func getMysqlInstance(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}

// Instance returns an instance of a MySQLRepo
func Instance(ctx context.Context, dns, dbName string, log logger.Logger) (*MysqlConn, *sql.DB, error) {
	return instance(ctx, dns, dbName, log, getMysqlInstance)
}

func instance(ctx context.Context, dsn, dbName string, log logger.Logger, getMysql getMysql) (*MysqlConn, *sql.DB, error) {
	var db *sql.DB
	onceDB.Do(func() {

		log.Info(ctx, "Connecting to database...")
		db, connErr = getMysql(dsn)
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
		if _, connErr = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)); connErr != nil {
			log.Errorf(ctx, "Create Database error: %v", connErr)
			return
		}

		if _, connErr = db.Exec(fmt.Sprintf("USE %s;", dbName)); connErr != nil {
			log.Errorf(ctx, "Default Database error: %v", connErr)
			return
		}

		connErr = mysqlDriver.SetLogger(log)
		if connErr != nil {
			return
		}
		log.Info(ctx, "Database successfully configured")

		conn = repoInstances(db)
		conn.db = db
	})

	return conn, db, connErr
}

func repoInstances(dbConn dbConn) *MysqlConn {
	return &MysqlConn{
		accountRepo: newAccountRepo(dbConn),
		authRepo:    newAuthRepo(dbConn),
	}
}

func (c *MysqlConn) WithTransaction(ctx context.Context, fn func(dm contract.DataManager) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	txConn := repoInstances(tx)
	err = fn(txConn)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return resterrors.NewInternalServerError("error rolling back transaction", rbErr.Error(), err.Error())
		}

		return err
	}

	return tx.Commit()
}

func (c *MysqlConn) DB() *sql.DB {
	return c.db
}

func (c *MysqlConn) Close() (err error) {
	return c.db.Close()
}

func (c *MysqlConn) Account() contract.AccountRepo {
	return c.accountRepo
}

func (c *MysqlConn) Auth() contract.AuthRepo {
	return c.authRepo
}
