package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"sync"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/resterrors"
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

// helps test the Instance function
type getMysql func(string) (*sql.DB, error)

func getMysqlInstance(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}

// Instance returns an instance of a MySQLRepo
func Instance(ctx context.Context, cfg *config.Config, log logger.Logger, migrationsDir string) (contract.DataManager, error) {
	return instance(ctx, cfg, log, migrationsDir, getMysqlInstance)
}

func instance(ctx context.Context, cfg *config.Config, log logger.Logger, migrationsDir string, getMysql getMysql) (contract.DataManager, error) {
	onceDB.Do(func() {
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=true",
			cfg.DB.MySQL.Username, cfg.DB.MySQL.Password, cfg.DB.MySQL.Host, cfg.DB.MySQL.Port,
		)

		var db *sql.DB
		log.Info(ctx, "Connecting to database...")
		db, connErr = getMysql(dataSourceName)
		if connErr != nil {
			return
		}

		cfg.AddCloser(func() {
			log.Info(ctx, "Closing mysql connection...")
			if err := db.Close(); err != nil {
				log.Errorf(ctx, "Error closing mysql connection: %v", err)
			}
		})

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
		connErr = migrate(db, migrationsDir)
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

func repoInstances(dbConn dbConn) *mysqlConn {
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
			return resterrors.NewInternalServerError("error rolling back transaction", rbErr.Error(), err.Error())
		}

		return err
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
