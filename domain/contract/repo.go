package contract

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/entity"
)

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	Begin() (Transaction, error)
	Account() AccountRepo
	Auth() AuthRepo
}

type Transaction interface {
	DataManager
	Rollback() error
	Commit() error
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session entity.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session entity.Session, err error)
}

// TODO: check if transfer should go to a transfer repo
type AccountRepo interface {
	AddTransfer(ctx context.Context, transferUUID string, accountOriginID, accountDestinationID int64, amount float64) (err error)
	CreateAccount(ctx context.Context, account entity.Account) (err error)
	GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
	GetTransfersByAccountID(ctx context.Context, accountID int64, origin bool) (transfers []entity.Transfer, err error)
	UpdateAccountBalance(ctx context.Context, accountID int64, balance float64) (err error)
}
