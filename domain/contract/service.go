package contract

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
)

type AccountService interface {
	CreateAccount(ctx context.Context, input dto.AccountInput) (err error)
	AddBalance(ctx context.Context, input dto.AddBalanceInput) (err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
	GetLoggedAccount(ctx context.Context) (account entity.Account, err error)
	GetLoggedAccountID(ctx context.Context) (accountID int64, err error)
}

type AuthService interface {
	Login(ctx context.Context, input dto.LoginInput) (account entity.Account, err error)
	CreateSession(ctx context.Context, session dto.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
	Logout(ctx context.Context, accessToken string) (err error)
}

type TransferService interface {
	CreateTransfer(ctx context.Context, transfer dto.TransferInput) (err error)
	GetTransfers(ctx context.Context, take, skip int64) (transfers []entity.Transfer, totalRecords int64, err error)
}
