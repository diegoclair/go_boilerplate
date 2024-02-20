package service

import (
	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
)

type Services struct {
	AccountService  contract.AccountService
	AuthService     contract.AuthService
	TransferService contract.TransferService
}

// New to get instance of all services
func New(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, crypto contract.Crypto, log logger.Logger) (*Services, error) {

	svc := newService(dm, cfg, cache, crypto, log)

	return &Services{
		AccountService:  newAccountService(svc),
		AuthService:     newAuthService(svc),
		TransferService: newTransferService(svc),
	}, nil
}

type service struct {
	dm     contract.DataManager
	cfg    *config.Config
	cache  contract.CacheManager
	log    logger.Logger
	crypto contract.Crypto
}

// newService has instances that will be used by the specific services
func newService(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, crypto contract.Crypto, log logger.Logger) *service {
	svc := new(service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cache = cache
	svc.log = log
	svc.crypto = crypto

	return svc
}
