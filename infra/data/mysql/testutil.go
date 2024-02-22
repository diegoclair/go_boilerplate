package mysql

import (
	"context"
	"log"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// SetMysqlTestContainerConfig set the mysql container for testing
//
// You can use this function to set the mysql container for an integration testing
func SetMysqlTestContainerConfig(ctx context.Context, cfg *config.Config) (closeFunc func()) {

	mysqlContainer, err := mysql.RunContainer(
		ctx,
		testcontainers.WithImage("mysql:8.0.32"),
		mysql.WithDatabase(cfg.DB.MySQL.DBName),
		mysql.WithUsername(cfg.DB.MySQL.Username),
		mysql.WithPassword(cfg.DB.MySQL.Password),
	)
	if err != nil {
		log.Fatalf("cannot start mysql container: %v", err)
	}

	cfg.DB.MySQL.Host, err = mysqlContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	cfg.DB.MySQL.Port = port.Port()

	return func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}
}
