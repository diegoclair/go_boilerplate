package service

import (
	"context"

	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/config"
)

type Service struct {
	dm     contract.Manager
	cfg    *config.Config
	cipher datacrypto.Crypto
}

func New(dm contract.Manager, cfg *config.Config, cipher datacrypto.Crypto) *Service {
	svc := new(Service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cipher = cipher

	return svc
}

type Manager interface {
	AccountService(svc *Service) AccountService
	AuthService(svc *Service) AuthService
	TransferService(svc *Service) TransferService
}

type PingService interface {
}

type AccountService interface {
	CreateAccount(account entity.Account) (err error)
	GetAccounts() (accounts []entity.Account, err error)
	GetAccountByUUID(accountUUID string) (account entity.Account, err error)
}

type AuthService interface {
	Login(cpf, secret string) (retVal entity.Authentication, err error)
}

type TransferService interface {
	CreateTransfer(appContext context.Context, transfer entity.Transfer) (err error)
	GetTransfers(appContext context.Context) (transfers []entity.Transfer, err error)
}

type serviceManager struct {
}

func NewServiceManager() Manager {
	return &serviceManager{}
}

func (s *serviceManager) AccountService(svc *Service) AccountService {
	return newAccountService(svc)
}

func (s *serviceManager) AuthService(svc *Service) AuthService {
	return newAuthService(svc)
}

func (s *serviceManager) TransferService(svc *Service) TransferService {
	return newTransferService(svc)
}
