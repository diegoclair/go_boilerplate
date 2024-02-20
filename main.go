package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	application "github.com/diegoclair/go_boilerplate/application/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/cache"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data"
	infraLogger "github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/transport/rest"
	"github.com/diegoclair/go_boilerplate/util/crypto"
	"github.com/diegoclair/go_utils-lib/v2/logger"
	"github.com/diegoclair/go_utils-lib/v2/validator"
)

const (
	gracefulShutdownTimeout = 10 * time.Second
)

func main() {

	cfg, err := config.GetConfigEnvironment(config.ProfileRun)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}
	defer cfg.Close()

	ctx := context.Background()
	log := infraLogger.New(cfg)

	authToken, err := auth.NewAuthToken(cfg.App.Auth, log)
	if err != nil {
		log.Fatalf(ctx, "Error getting NewAuthToken: %v", err)
	}

	data, err := data.Connect(ctx, cfg, log)
	if err != nil {
		log.Fatalf(ctx, "Error to connect dataManager repositories: %v", err)
	}

	log.Infof(ctx, "Connecting to the cache server at %s:%d.", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port)
	cache, err := cache.Instance(ctx, cfg, log)
	if err != nil {
		log.Fatalf(ctx, "Error connecting to cache server: %v", err)
	}

	c := crypto.NewCrypto()

	apps, err := application.New(data, cfg, cache, c, log)
	if err != nil {
		log.Fatalf(ctx, "error to get domain services: %v", err)
	}

	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf(ctx, "error to get validator: %v", err)
	}

	server := rest.StartRestServer(ctx, cfg, apps, log, authToken, v)

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

	if err := server.Srv.Shutdown(ctx); err != nil {
		log.Errorf(ctx, "Error to shutdown rest server: %v", err)
	}
}
