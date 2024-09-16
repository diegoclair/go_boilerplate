package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diegoclair/go_boilerplate/application/service"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/migrator/mysql"
	"github.com/diegoclair/go_boilerplate/transport/rest"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
)

const (
	gracefulShutdownTimeout = 10 * time.Second
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfigEnvironment(ctx, "boilerplate")
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}
	defer cfg.Close()

	log := cfg.GetLogger()
	authToken := cfg.GetAuthToken()
	cache := cfg.GetCacheManager()
	c := cfg.GetCrypto()

	v, err := validator.NewValidator()
	if err != nil {
		log.Errorf(ctx, "error to get validator: %v", err)
		return
	}
	data := cfg.GetDataManager()

	log.Info(ctx, "Running the migrations...")
	err = mysql.Migrate(data.DB())
	if err != nil {
		log.Errorf(ctx, "error to migrate mysql: %v", err)
		return
	}
	log.Info(ctx, "Migrations completed successfully")

	services, err := service.New(
		service.WithDataManager(data),
		service.WithConfig(cfg),
		service.WithCacheManager(cache),
		service.WithLogger(log),
		service.WithCrypto(c),
		service.WithValidator(v),
	)
	if err != nil {
		log.Errorf(ctx, "error to get domain services: %v", err)
		return
	}

	server := rest.StartRestServer(ctx, cfg, services, log, authToken, cache)

	gracefulShutdown(server, log)
}

// will wait for a SIGTERM or SIGINT signal and wait the server to finish processing requests or timeout after 10 seconds
func gracefulShutdown(server *rest.Server, log logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	log.Info(ctx, "Shutting down...")

	if err := server.Router.Echo().Shutdown(ctx); err != nil {
		log.Errorf(ctx, "Error to shutdown rest server: %v", err)
	}
}
