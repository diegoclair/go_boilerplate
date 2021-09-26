package repo

import "github.com/diegoclair/go-boilerplate/domain/entity"

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

type AccountRepo interface {
	AddTransfer(transfer entity.Transfer) (err error)
	CreateAccount(account entity.Account) (err error)
	GetAccountByDocument(encryptedCPF string) (account entity.Account, err error)
	GetAccounts() (accounts []entity.Account, err error)
	GetAccountByUUID(accountUUID string) (account entity.Account, err error)
	GetTransfersByAccountID(accountID int64, origin bool) (transfers []entity.Transfer, err error)
	UpdateAccountBalance(account entity.Account) (err error)
}
