package mysql

import (
	"database/sql"

	"github.com/diegoclair/go-boilerplate/domain/contract"
)

type mysqlTransaction struct {
	tx         *sql.Tx
	committed  bool
	rolledback bool
}

func newTransaction(tx *sql.Tx) *mysqlTransaction {
	instance := &mysqlTransaction{tx: tx}
	return instance
}

func (t *mysqlTransaction) Begin() (contract.Transaction, error) {
	return &mysqlTransaction{
		tx: t.tx,
	}, nil
}

// Commit persists changes to database
func (t *mysqlTransaction) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		return err
	}

	t.committed = true

	return nil
}

// Rollback discards changes on database
func (t *mysqlTransaction) Rollback() error {
	if t != nil && !t.committed && !t.rolledback {
		err := t.tx.Rollback()
		if err != nil {
			return err
		}

		t.rolledback = true
	}

	return nil
}

func (t *mysqlTransaction) Account() contract.AccountRepo {
	return newAccountRepo(t.tx)
}
