package data

import (
	"github.com/diegoclair/go-boilerplate/domain/repo"
	"github.com/diegoclair/go-boilerplate/infra/data/mysql"
)

// Connect returns a instace of mysql db
func Connect() (repo.Manager, error) {
	return mysql.Instance()
}
