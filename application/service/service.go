package service

import (
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
)

type Services struct {
	AccountService  contract.AccountService
	AuthService     contract.AuthService
	TransferService contract.TransferService
}

type ServiceOptions func(*service)

// New to get instance of all services
func New(svcOptions ...ServiceOptions) (*Services, error) {

	svc := &service{}
	for _, opt := range svcOptions {
		opt(svc)
	}

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

func WithDataManager(dm contract.DataManager) ServiceOptions {
	return func(s *service) {
		s.dm = dm
	}
}

func WithConfig(cfg *config.Config) ServiceOptions {
	return func(s *service) {
		s.cfg = cfg
	}
}

func WithCacheManager(cache contract.CacheManager) ServiceOptions {
	return func(s *service) {
		s.cache = cache
	}
}

func WithLogger(log logger.Logger) ServiceOptions {
	return func(s *service) {
		s.log = log
	}
}

func WithCrypto(crypto contract.Crypto) ServiceOptions {
	return func(s *service) {
		s.crypto = crypto
	}
}

func WithValidator(v validator.Validator) ServiceOptions {
	return func(s *service) {
		s.validator = v
	}
}
