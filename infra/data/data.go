package data

import (
	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/infra/data/mysql"
	"github.com/diegoclair/go-boilerplate/util/config"
)

// Connect returns a instace of mysql db
func Connect(cfg *config.Config) (contract.Manager, error) {
	return mysql.Instance(cfg)
}
