package service

import (
	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
)

type Services struct {
	AccountService  contract.AccountService
	AuthService     contract.AuthService
	TransferService contract.TransferService
}

// New to get instance of all services
func New(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, crypto contract.Crypto, log logger.Logger, v validator.Validator) (*Services, error) {

	svc := newService(dm, cfg, cache, crypto, log, v)

	return &Services{
		AccountService:  newAccountService(svc),
		AuthService:     newAuthService(svc),
		TransferService: newTransferService(svc),
	}, nil
}

type service struct {
	dm        contract.DataManager
	cfg       *config.Config
	cache     contract.CacheManager
	log       logger.Logger
	crypto    contract.Crypto
	validator validator.Validator
}

// newService has instances that will be used by the specific services
func newService(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, crypto contract.Crypto, log logger.Logger, validator validator.Validator) *service {
	svc := new(service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cache = cache
	svc.log = log
	svc.crypto = crypto
	svc.validator = validator

	return svc
}
