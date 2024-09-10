package mysql

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/configmock"
)

var (
	testMysql contract.DataManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg := configmock.New()

	close := setMysqlTestContainerConfig(ctx, cfg)
	defer close()

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting root dir: %v", err)
	}

	migrationsDir := rootDir + "/../migrations/mysql"

	mysql, _, err := Instance(ctx,
		cfg.DB.MySQL.Host,
		cfg.DB.MySQL.Port,
		cfg.DB.MySQL.Username,
		cfg.DB.MySQL.Password,
		cfg.DB.MySQL.DBName,
		cfg.GetLogger(),
		migrationsDir,
	)
	if err != nil {
		log.Fatalf("cannot connect to mysql: %v", err)
	}

	testMysql = mysql

	os.Exit(m.Run())
}
