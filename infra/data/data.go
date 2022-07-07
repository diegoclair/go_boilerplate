package data

import (
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/data/mysql"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/util/config"
)

// Connect returns a instace of mysql db
func Connect(cfg *config.Config, log logger.Logger) (contract.DataManager, error) {
	return mysql.Instance(cfg, log)
}
