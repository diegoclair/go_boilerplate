package service

import (
	"errors"
	"time"

	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
)

type Apps struct {
	AccountService  contract.AccountApp
	AuthService     contract.AuthApp
	TransferService contract.TransferApp
}

// New to get instance of all services
func New(infra domain.Infrastructure, accessTokenDuration time.Duration) (*Apps, error) {
	if err := validateInfrastructure(infra); err != nil {
		return nil, err
	}

	accSvc := newAccountService(infra)

	return &Apps{
		AccountService:  accSvc,
		AuthService:     newAuthApp(infra, accSvc, accessTokenDuration),
		TransferService: newTransferService(infra, accSvc),
	}, nil
}

// validateInfrastructure validate the dependencies needed to initialize the services
func validateInfrastructure(infra domain.Infrastructure) error {
	if infra.Logger() == nil {
		return errors.New("logger is required")
	}

	if infra.DataManager() == nil {
		return errors.New("data manager is required")
	}

	if infra.CacheManager() == nil {
		return errors.New("cache manager is required")
	}

	if infra.Crypto() == nil {
		return errors.New("crypto is required")
	}

	if infra.Validator() == nil {
		return errors.New("validator is required")
	}

	return nil
}
