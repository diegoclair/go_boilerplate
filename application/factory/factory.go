package factory

import (
	"log"
	"sync"

	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/domain/service"
	"github.com/diegoclair/go-boilerplate/infra/data"
	"github.com/diegoclair/go-boilerplate/util/config"
)

type services struct {
	Mapper          mapper.Mapper
	AccountService  service.AccountService
	AuthService     service.AuthService
	TransferService service.TransferService
}

var (
	instance *services
	once     sync.Once
)

//GetDomainServices to get instace of all services
func GetDomainServices(cfg *config.Config) *services {

	once.Do(func() {

		data, err := data.Connect(cfg)
		if err != nil {
			log.Fatalf("Error to connect data repositories: %v", err)
		}

		instance = &services{}
		cipher := datacrypto.NewAESECB(datacrypto.AES256, cfg.DB.MySQL.CryptoKey)
		svc := service.New(data, cfg, cipher)
		svm := service.NewServiceManager()

		instance.Mapper = mapper.New()
		instance.AccountService = svm.AccountService(svc)
		instance.AuthService = svm.AuthService(svc)
		instance.TransferService = svm.TransferService(svc)
	})

	return instance
}
