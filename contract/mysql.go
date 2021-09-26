package contract

import "github.com/diegoclair/go-boilerplate/domain/entity"

// MysqlTransaction holds the methods that manipulates the main data, from within a transaction.
type MysqlTransaction interface {
	MySQLRepo
	Rollback() error
	Commit() error
}

//MySQLRepo defines the repository aggregator interface
type MySQLRepo interface {
	Begin() (MysqlTransaction, error)
	Account() AccountRepo
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
