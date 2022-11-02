package mysql

import (
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

		conn = &mysqlConn{
			db: db,
		}
	})

	return conn, connErr
}

// Begin starts a mysql transaction
func (c *mysqlConn) Begin() (contract.Transaction, error) {
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

func (c *mysqlConn) Auth() contract.AuthRepo {
	return newAuthRepo(c.db)
}
