package service

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
)

type Services struct {
	AccountService  AccountService
	AuthService     AuthService
	TransferService TransferService
}

// New to get instace of all services
func New(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, log logger.Logger) (*Services, error) {

	svc := newService(dm, cfg, cache, log)

	return &Services{
		AccountService:  newAccountService(svc),
		AuthService:     newAuthService(svc),
		TransferService: newTransferService(svc),
	}, nil
}

type service struct {
	dm    contract.DataManager
	cfg   *config.Config
	cache contract.CacheManager
	log   logger.Logger
}

// newService has instances that will be used by the specific services
func newService(dm contract.DataManager, cfg *config.Config, cache contract.CacheManager, log logger.Logger) *service {
	svc := new(service)
	svc.dm = dm
	svc.cfg = cfg
	svc.cache = cache
	svc.log = log

	return svc
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
