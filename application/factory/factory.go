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

type Services struct {
	Cfg             *config.Config
	Mapper          mapper.Mapper
	AccountService  service.AccountService
	AuthService     service.AuthService
	TransferService service.TransferService
}

var (
	instance *Services
	once     sync.Once
)

//GetDomainServices to get instace of all services
func GetDomainServices() *Services {

	once.Do(func() {

		data, err := data.Connect()
		if err != nil {
			log.Fatalf("Error to connect data repositories: %v", err)
		}

		instance = &Services{}
		cfg := config.GetConfigEnvironment()
		cipher := datacrypto.NewAESECB(datacrypto.AES256, cfg.MySQL.CryptoKey)
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
