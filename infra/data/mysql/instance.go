package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"sync"

	"github.com/GuiaBolso/darwin"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data/migrations"
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
func Instance(cfg *config.Config, log logger.Logger) (contract.DataManager, error) {
	onceDB.Do(func() {

		dataSourceName := fmt.Sprintf("%s:root@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
			cfg.DB.MySQL.Username, cfg.DB.MySQL.Host, cfg.DB.MySQL.Port, cfg.DB.MySQL.DBName,
		)

		var db *sql.DB
		log.Info("Connecting to database...")
		db, connErr = sql.Open("mysql", dataSourceName)
		if connErr != nil {
			return
		}

		log.Info("Database Ping...")
		connErr = db.Ping()
		if connErr != nil {
			return
		}

		log.Info("Creating database if not exists...")
		if _, connErr = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.DB.MySQL.DBName)); connErr != nil {
			log.Error("Create Database error: ", connErr)
			return
		}

		if _, connErr = db.Exec(fmt.Sprintf("USE %s;", cfg.DB.MySQL.DBName)); connErr != nil {
			log.Error("Default Database error: ", connErr)
			return
		}

		connErr = mysqlDriver.SetLogger(log)
		if connErr != nil {
			return
		}
		log.Info("Database successfully configured")

		log.Info("Running the migrations")
		driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})

		d := darwin.New(driver, migrations.Migrations, nil)

		connErr = d.Migrate()
		if connErr != nil {
			log.Error("Migrate Error: ", connErr)
			return
		}

		log.Info("Migrations executed")

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
