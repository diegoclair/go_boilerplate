package contract

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/domain/transfer"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account account.Account) (err error)
	AddBalance(ctx context.Context, accountUUID string, amount float64) (err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []account.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account account.Account, err error)
}

type AuthService interface {
	Login(ctx context.Context, cpf, secret string) (account account.Account, err error)
	CreateSession(ctx context.Context, session dto.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
}

type TransferService interface {
	CreateTransfer(ctx context.Context, transfer transfer.Transfer) (err error)
	GetTransfers(ctx context.Context) (transfers []transfer.Transfer, err error)
}
