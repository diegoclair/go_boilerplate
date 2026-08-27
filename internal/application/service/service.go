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

func New(infra domain.Infrastructure, accessTokenDuration time.Duration) (*Apps, error) {
	if err := validateInfrastructure(infra); err != nil {
		return nil, err
	}

	accSvc := newAccountService(infra)

	authSvc, err := newAuthApp(infra, accSvc, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &Apps{
		AccountService:  accSvc,
		AuthService:     authSvc,
		TransferService: newTransferService(infra, accSvc),
	}, nil
}

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
