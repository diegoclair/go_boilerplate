/*
Copyright © 2021 Diego Clair Rodrigues <diego93rodrigues@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diegoclair/go_boilerplate/application/rest"
	"github.com/diegoclair/go_boilerplate/domain/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/cache"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/util/crypto"
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
	log := logger.New(cfg)

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

	services, err := service.New(data, cfg, cache, c, log)
	if err != nil {
		log.Fatalf(ctx, "error to get domain services: %v", err)
	}

	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf(ctx, "error to get validator: %v", err)
	}

	server := rest.StartRestServer(ctx, cfg, services, log, authToken, v)

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
