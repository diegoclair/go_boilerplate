package factory

import (
	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/domain/service"
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
func GetServices(cfg *config.Config, data contract.DataManager, svc *service.Service, svm service.Manager, log logger.Logger, cipher datacrypto.Crypto) (*Services, error) {

	services := &Services{}

	services.Mapper = mapper.New()
	services.AccountService = svm.AccountService(svc)
	services.AuthService = svm.AuthService(svc)
	services.TransferService = svm.TransferService(svc)

	return services, nil
}
