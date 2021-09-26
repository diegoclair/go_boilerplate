package mysql

import (
	"database/sql"
	"log"

	"github.com/diegoclair/go-boilerplate/contract"
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

func (t *mysqlTransaction) Begin() (contract.MysqlTransaction, error) {
	return &mysqlTransaction{
		tx: t.tx,
	}, nil
}

func (t *mysqlTransaction) MySQL() contract.MySQLRepo {
	mysqlRepo, err := Instance()
	if err != nil {
		log.Fatalf("Error to start mysql instance: %v", err)
	}
	return mysqlRepo
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
