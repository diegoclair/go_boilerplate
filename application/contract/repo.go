package contract

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
)

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	WithTransaction(ctx context.Context, fn func(r DataManager) error) error
	Account() AccountRepo
	Auth() AuthRepo
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session dto.Session) (err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
}

type AccountRepo interface {
	AddTransfer(ctx context.Context, transferUUID string, accountOriginID, accountDestinationID int64, amount float64) (err error)
	CreateAccount(ctx context.Context, account entity.Account) (createdID int64, err error)
	GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
	GetTransfersByAccountID(ctx context.Context, accountID, take, skip int64, origin bool) (transfers []entity.Transfer, totalRecords int64, err error)
	UpdateAccountBalance(ctx context.Context, accountID int64, balance float64) (err error)
}
