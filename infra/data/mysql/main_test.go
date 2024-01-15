package mysql_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data/mysql"
	"github.com/diegoclair/go_boilerplate/infra/logger"
)

var testMysql contract.DataManager

func TestMain(m *testing.M) {

	cfg, err := config.GetConfigEnvironment("./../../../config.toml")
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting root dir: %v", err)
	}
	migrationsDir := rootDir + "/../migrations/mysql"

	log := logger.NewNoop()
	cfg.DB.MySQL.DBName = cfg.DB.MySQL.DBName + "_test"

	ctx := context.Background()
	mysql, err := mysql.Instance(ctx, cfg, log, migrationsDir)
	if err != nil {
		log.Fatalf(ctx, "cannot connect to mysql: %v", err)
	}

	testMysql = mysql
	os.Exit(m.Run())
}
