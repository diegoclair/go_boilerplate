package domain

import (
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/validator"
)

type Infrastructure interface {
	CacheManager() contract.CacheManager
	DataManager() contract.DataManager
	Logger() logger.Logger
	Crypto() contract.Crypto
	Validator() validator.Validator
}

type infrastructureServices struct {
	cacheManager contract.CacheManager
	dataManager  contract.DataManager
	logger       logger.Logger
	crypto       contract.Crypto
	validator    validator.Validator
}

type InfraOption func(*infrastructureServices)

func WithCacheManager(cacheManager contract.CacheManager) InfraOption {
	return func(i *infrastructureServices) {
		i.cacheManager = cacheManager
	}
}

func WithDataManager(dataManager contract.DataManager) InfraOption {
	return func(i *infrastructureServices) {
		i.dataManager = dataManager
	}
}

func WithLogger(logger logger.Logger) InfraOption {
	return func(i *infrastructureServices) {
		i.logger = logger
	}
}

func WithCrypto(crypto contract.Crypto) InfraOption {
	return func(i *infrastructureServices) {
		i.crypto = crypto
	}
}

func WithValidator(validator validator.Validator) InfraOption {
	return func(i *infrastructureServices) {
		i.validator = validator
	}
}

func NewInfrastructureServices(options ...InfraOption) Infrastructure {
	infra := &infrastructureServices{}
	for _, option := range options {
		option(infra)
	}
	return infra
}

func (i *infrastructureServices) CacheManager() contract.CacheManager {
	return i.cacheManager
}

func (i *infrastructureServices) DataManager() contract.DataManager {
	return i.dataManager
}

func (i *infrastructureServices) Logger() logger.Logger {
	return i.logger
}

func (i *infrastructureServices) Crypto() contract.Crypto {
	return i.crypto
}

func (i *infrastructureServices) Validator() validator.Validator {
	return i.validator
}
