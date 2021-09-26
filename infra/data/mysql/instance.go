package mysql

import (
	"database/sql"
	"fmt"

	"sync"

	"github.com/GuiaBolso/darwin"
	"github.com/labstack/gommon/log"

	"github.com/diegoclair/go-boilerplate/contract"
	"github.com/diegoclair/go-boilerplate/infra/data/migrations"
	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/diegoclair/go_utils-lib/logger"
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
}

//Instance returns an instance of a MySQLRepo
func Instance() (contract.MySQLRepo, error) {
	onceDB.Do(func() {
		cfg := config.GetConfigEnvironment()

		dataSourceName := fmt.Sprintf("%s:root@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
			cfg.MySQL.Username, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.DBName,
		)

		log.Info("Connecting to database...")
		db, connErr := sql.Open("mysql", dataSourceName)
		if connErr != nil {
			return
		}

		log.Info("Database Ping...")
		connErr = db.Ping()
		if connErr != nil {
			return
		}

		log.Info("Creating database...")
		if _, connErr = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.MySQL.DBName)); connErr != nil {
			logger.Error("Create Database error: ", connErr)
			return
		}

		if _, connErr = db.Exec(fmt.Sprintf("USE %s;", cfg.MySQL.DBName)); connErr != nil {
			logger.Error("Default Database error: ", connErr)
			return
		}

		connErr = mysqlDriver.SetLogger(logger.GetLogger())
		if connErr != nil {
			return
		}
		logger.Info("Database successfully configured")

		logger.Info("Running the migrations")
		driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})

		d := darwin.New(driver, migrations.Migrations, nil)

		connErr = d.Migrate()
		if connErr != nil {
			logger.Error("Migrate Error: ", connErr)
			return
		}

		logger.Info("Migrations executed")

		conn = &mysqlConn{
			db: db,
		}
	})

	return conn, connErr
}

// Begin starts a transaction
func (c *mysqlConn) Begin() (contract.MysqlTransaction, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	return newTransaction(tx), nil
}

func (c *mysqlConn) Close() (err error) {
	return c.db.Close()
}

func (c *mysqlConn) Account() contract.AccountRepo {
	return newAccountRepo(c.db)
}
