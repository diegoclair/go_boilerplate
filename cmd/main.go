package main

import (
	"context"
	"log"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	db "github.com/diegoclair/go_boilerplate/infra/data/mysql"
	"github.com/diegoclair/go_boilerplate/infra/shutdown"
	"github.com/diegoclair/go_boilerplate/internal/application/service"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest"
	"github.com/diegoclair/go_boilerplate/migrator/mysql"
	"github.com/diegoclair/go_utils/logger"
)

const (
	gracefulShutdownTimeout = 10 * time.Second
	appName                 = "boilerplate"
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfigEnvironment(ctx, appName)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}
	defer cfg.Close()

	log := cfg.GetLogger()

	infra := domain.NewInfrastructureServices(
		domain.WithCacheManager(cfg.GetCacheManager()),
		domain.WithDataManager(cfg.GetDataManager()),
		domain.WithLogger(log),
		domain.WithCrypto(cfg.GetCrypto()),
		domain.WithValidator(cfg.GetValidator()),
	)

	log.Info(ctx, "Running the migrations...")
	err = mysql.Migrate(cfg.GetDataManager().(*db.MysqlConn).DB())
	if err != nil {
		log.Errorw(ctx, "error to migrate mysql", logger.Err(err))
		return
	}
	log.Info(ctx, "Migrations completed successfully")

	apps, err := service.New(infra, cfg.App.Auth.AccessTokenDuration)
	if err != nil {
		log.Errorw(ctx, "error to get domain services", logger.Err(err))
		return
	}

	server := rest.StartRestServer(ctx, cfg, infra, apps, appName, cfg.GetHttpPort())

	shutdown.GracefulShutdown(ctx, log, shutdown.WithRestServer(server.Router.Echo()))
}
