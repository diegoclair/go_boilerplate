package mysql

import (
	"log"
	"os"
	"testing"

	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/util/config"
)

var testMysql contract.Manager

func TestMain(m *testing.M) {

	cfg, err := config.GetConfigEnvironment("./../../../config.toml")
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}

	mysql, err := Instance(cfg)
	if err != nil {
		log.Fatal("cannot connect to mysql: ", err)
	}

	testMysql = mysql
	os.Exit(m.Run())
}
