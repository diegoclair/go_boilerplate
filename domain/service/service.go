package service

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/util/config"
)

type Service struct {
	dm    contract.DataManager
	cfg   *config.Config
	cache contract.CacheManager
	log   logger.Logger
}

func New(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, log logger.Logger) *Service {
	svc := new(Service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cache = cache
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
	AddBalance(ctx context.Context, accountUUID string, amount float64) (err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
}

type AuthService interface {
	Login(ctx context.Context, cpf, secret string) (account entity.Account, err error)
	CreateSession(ctx context.Context, session entity.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session entity.Session, err error)
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
