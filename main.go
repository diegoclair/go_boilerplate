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

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest"
	"github.com/diegoclair/go_boilerplate/domain/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/cache"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/data"
	"github.com/diegoclair/go_boilerplate/infra/logger"
)

func main() {

	cfg, err := config.GetConfigEnvironment(config.ConfigDefaultName)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}

	ctx := context.Background()
	log := logger.New(*cfg)

	authToken, err := auth.NewAuthToken(cfg.App.Auth)
	if err != nil {
		log.Fatalf(ctx, "Error to load config: %v", err)
	}

	data, err := data.Connect(ctx, cfg, log)
	if err != nil {
		log.Fatalf(ctx, "Error to connect dataManager repositories: %v", err)
	}

	log.Infof(ctx, "Connecting to the cache server at %s:%d.", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port)
	cache, err := cache.Instance(cfg.Cache.Redis, log)
	if err != nil {
		log.Fatalf(ctx, "Error connecting to cache server: %v", err)
	}

	services, err := service.New(data, cfg, cache, log)
	if err != nil {
		log.Fatalf(ctx, "error to get domain services: %v", err)
	}

	mp := mapper.New()

	rest.StartRestServer(ctx, cfg, services, log, authToken, mp) //TODO: receive flags for what server it will starts
}
