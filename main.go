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

	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/diegoclair/go-boilerplate/application/factory"
	"github.com/diegoclair/go-boilerplate/application/rest"
	"github.com/diegoclair/go-boilerplate/domain/service"
	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go-boilerplate/infra/cache"
	"github.com/diegoclair/go-boilerplate/infra/data"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"
)

func main() {

	cfg, err := config.GetConfigEnvironment(config.ConfigDefaultFilepath)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}

	log := logger.New(cfg.Log, cfg.App.Name)
	cipher := datacrypto.NewAESECB(datacrypto.AES256, cfg.DB.MySQL.CryptoKey)

	authToken, err := auth.NewAuthToken(cfg.App.Auth)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}

	data, err := data.Connect(cfg, log)
	if err != nil {
		log.Fatalf("Error to connect dataManager repositories: %v", err)
	}

	log.Info("Connecting to the cache server at %s:%d.", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port)
	cache, err := cache.Instance(cfg.Cache.Redis, log)
	if err != nil {
		log.Fatal("Error connecting to cache server:", err)
	}

	svc := service.New(data, cfg, cache, cipher, log)
	svm := service.NewServiceManager()

	services, err := factory.GetServices(svc, svm)
	if err != nil {
		log.Fatal("error to get domain services: ", err)
	}

	rest.StartRestServer(cfg, services, log, authToken) //TODO: receive flags for what server it will starts
}
