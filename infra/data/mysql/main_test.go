package mysql

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_boilerplate/migrator/mysql"
)

var (
	testMysql contract.DataManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg := configmock.New()

	close := setMysqlTestContainerConfig(ctx, cfg)
	defer close()

	mysqlConn, db, err := Instance(ctx,
		cfg.GetMysqlDNS(),
		cfg.DB.MySQL.DBName,
		cfg.GetLogger(),
	)
	if err != nil {
		log.Fatalf("cannot connect to mysql: %v", err)
	}

	err = mysql.Migrate(db)
	if err != nil {
		log.Fatalf("cannot migrate mysql: %v", err)
	}

	testMysql = mysqlConn

	os.Exit(m.Run())
}
