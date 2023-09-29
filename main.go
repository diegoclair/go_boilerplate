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

func logsome(cfg *config.Config) {

}

func main() {

	cfg, err := config.GetConfigEnvironment(config.ConfigDefaultName)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}

	log := logger.New(*cfg)

	log.Info("info ", "message4", "message2", "message3", "message1")
	log.Fatal("fatal ", "message1", "message2")

	authToken, err := auth.NewAuthToken(cfg.App.Auth)
	if err != nil {
		log.Fatal("Error to load config: ", err)
	}

	data, err := data.Connect(cfg, log)
	if err != nil {
		log.Fatal("Error to connect dataManager repositories: ", err)
	}

	log.Info("Connecting to the cache server at %s:%d.", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port)
	cache, err := cache.Instance(cfg.Cache.Redis, log)
	if err != nil {
		log.Fatal("Error connecting to cache server:", err)
	}

	services, err := service.New(data, cfg, cache, log)
	if err != nil {
		log.Fatal("error to get domain services: ", err)
	}

	mp := mapper.New()

	rest.StartRestServer(cfg, services, log, authToken, mp) //TODO: receive flags for what server it will starts
}
