package contract

import (
	"context"

	"github.com/diegoclair/go-boilerplate/domain/entity"
)

// Manager holds the methods that manipulates the main data.
type Manager interface {
	Begin() (Transaction, error)
	Account() AccountRepo
}

type Transaction interface {
	Manager
	Rollback() error
	Commit() error
}

//TODO: check if transfer should go to a transfer repo
type AccountRepo interface {
	AddTransfer(ctx context.Context, transfer entity.Transfer) (err error)
	CreateAccount(ctx context.Context, account entity.Account) (err error)
	GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error)
	GetAccounts(ctx context.Context) (accounts []entity.Account, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
	GetTransfersByAccountID(ctx context.Context, accountID int64, origin bool) (transfers []entity.Transfer, err error)
	UpdateAccountBalance(ctx context.Context, account entity.Account) (err error)
}
