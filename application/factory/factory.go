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
	Cfg             *config.Config
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
func GetDomainServices() *services {

	once.Do(func() {
		cfg, err := config.GetConfigEnvironment()
		if err != nil {
			log.Fatalf("Error to load config: %v", err)
		}

		data, err := data.Connect(cfg)
		if err != nil {
			log.Fatalf("Error to connect data repositories: %v", err)
		}

		instance = &services{}
		cipher := datacrypto.NewAESECB(datacrypto.AES256, cfg.DB.MySQL.CryptoKey)
		svc := service.New(data, cfg, cipher)
		svm := service.NewServiceManager()

		instance.Cfg = cfg
		instance.Mapper = mapper.New()
		instance.AccountService = svm.AccountService(svc)
		instance.AuthService = svm.AuthService(svc)
		instance.TransferService = svm.TransferService(svc)
	})

	return instance
}
