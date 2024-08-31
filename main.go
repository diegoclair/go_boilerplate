package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diegoclair/go_boilerplate/application/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/cache"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data"
	infraLogger "github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/transport/rest"
	"github.com/diegoclair/go_boilerplate/util/crypto"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
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
		log.Errorf(ctx, "Error getting NewAuthToken: %v", err)
		return
	}

	data, err := data.Connect(ctx, cfg, log)
	if err != nil {
		log.Errorf(ctx, "Error to connect dataManager repositories: %v", err)
		return
	}

	log.Infof(ctx, "Connecting to the cache server at %s:%d.", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port)
	cache, err := cache.Instance(ctx, cfg, log)
	if err != nil {
		log.Errorf(ctx, "Error connecting to cache server: %v", err)
		return
	}

	c := crypto.NewCrypto()

	v, err := validator.NewValidator()
	if err != nil {
		log.Errorf(ctx, "error to get validator: %v", err)
		return
	}

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
