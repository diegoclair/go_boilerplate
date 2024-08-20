package mysql

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
)

var (
	testMysql contract.DataManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg, err := config.GetConfigEnvironment(config.ProfileTest)
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	close := SetMysqlTestContainerConfig(ctx, cfg)
	defer close()

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting root dir: %v", err)
	}

	migrationsDir := rootDir + "/../migrations/mysql"

	mysql, err := Instance(ctx, cfg, logger.NewNoop(), migrationsDir)
	if err != nil {
		log.Fatalf("cannot connect to mysql: %v", err)
	}

	testMysql = mysql

	os.Exit(m.Run())
}
