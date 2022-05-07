package service

import (
	"context"

	"github.com/IQ-tech/go-crypto-layer/datacrypto"
	"github.com/diegoclair/go-boilerplate/domain/contract"
	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"
)

type Service struct {
	dm     contract.Manager
	cfg    *config.Config
	cipher datacrypto.Crypto
	log    logger.Logger
}

func New(dm contract.Manager, cfg *config.Config, cipher datacrypto.Crypto, log logger.Logger) *Service {
	svc := new(Service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cipher = cipher
	svc.log = log

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
	CreateAccount(ctx context.Context, account entity.Account) (err error)
	GetAccounts(ctx context.Context) (accounts []entity.Account, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
}

type AuthService interface {
	Login(ctx context.Context, cpf, secret string) (account entity.Account, err error)
}

type TransferService interface {
	CreateTransfer(ctx context.Context, transfer entity.Transfer) (err error)
	GetTransfers(ctx context.Context) (transfers []entity.Transfer, err error)
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
