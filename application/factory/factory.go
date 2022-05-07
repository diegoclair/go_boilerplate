package factory

import (
	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/domain/service"
	"github.com/diegoclair/go-boilerplate/infra/data"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"
)

type Services struct {
	Mapper          mapper.Mapper
	AccountService  service.AccountService
	AuthService     service.AuthService
	TransferService service.TransferService
}

//GetServices to get instace of all services
func GetServices(cfg *config.Config, log logger.Logger, cipher datacrypto.Crypto) (*Services, error) {

	data, err := data.Connect(cfg, log)
	if err != nil {
		log.Error("Error to connect data repositories: %v", err)
		return nil, err
	}

	services := &Services{}
	svc := service.New(data, cfg, cipher, log)
	svm := service.NewServiceManager()

	services.Mapper = mapper.New()
	services.AccountService = svm.AccountService(svc)
	services.AuthService = svm.AuthService(svc)
	services.TransferService = svm.TransferService(svc)

	return services, nil
}
