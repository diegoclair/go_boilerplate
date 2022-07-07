package factory

import (
	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/domain/service"
)

type Services struct {
	Mapper          mapper.Mapper
	AccountService  service.AccountService
	AuthService     service.AuthService
	TransferService service.TransferService
}

//GetServices to get instace of all services
func GetServices(svc *service.Service, svm service.Manager) (*Services, error) {

	services := &Services{}

	services.Mapper = mapper.New()
	services.AccountService = svm.AccountService(svc)
	services.AuthService = svm.AuthService(svc)
	services.TransferService = svm.TransferService(svc)

	return services, nil
}
