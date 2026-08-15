package contract

import (
	"context"

	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
)

// Repos is every repository and nothing else. A transaction hands this to its
// callback, so opening a second one from inside cannot compile.
type Repos interface {
	Account() AccountRepo
	Auth() AuthRepo
}

// DataManager holds the methods that manipulates the main data.
type DataManager interface {
	Repos
	WithTransaction(ctx context.Context, fn func(tx Repos) error) error
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session dto.Session) (sessionID int64, err error)
	GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error)
	SetSessionAsBlocked(ctx context.Context, sessionUUID string) (err error)
}

type AccountRepo interface {
	AddTransfer(ctx context.Context, transferUUID string, accountOriginID, accountDestinationID int64, amount float64) (transferID int64, err error)
	CreateAccount(ctx context.Context, account entity.Account) (createdID int64, err error)
	GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error)
	GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error)
	GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error)
	GetAccountIDByUUID(ctx context.Context, accountUUID string) (accountID int64, err error)
	GetTransfersByAccountID(ctx context.Context, accountID, take, skip int64, origin bool) (transfers []entity.Transfer, totalRecords int64, err error)
	UpdateAccountBalance(ctx context.Context, accountID int64, balance float64) (err error)
}
