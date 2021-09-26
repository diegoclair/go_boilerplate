package data

import (
	"log"

	"github.com/diegoclair/go-boilerplate/contract"
	"github.com/diegoclair/go-boilerplate/infra/data/mysql"
)

//we can add here more than one database
type data struct {
	mysqlRepo contract.MySQLRepo
}

// Connect returns a instace of mysql db
func Connect() (contract.DataManager, error) {
	repo := new(data)
	return &data{
		mysqlRepo: repo.MySQL(),
	}, nil
}

func (d *data) MySQL() contract.MySQLRepo {
	mysqlRepo, err := mysql.Instance()
	if err != nil {
		log.Fatalf("Error to start mysql instance: %v", err)
	}
	return mysqlRepo
}
