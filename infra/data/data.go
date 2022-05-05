package data

import (
	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/infra/data/mysql"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"
)

// Connect returns a instace of mysql db
func Connect(cfg *config.Config, log logger.Logger) (contract.Manager, error) {
	return mysql.Instance(cfg, log)
}
