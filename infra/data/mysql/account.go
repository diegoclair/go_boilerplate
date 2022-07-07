package mysql

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils-lib/v2/mysqlutils"
)

type accountRepo struct {
	db connection
}

func newAccountRepo(db connection) contract.AccountRepo {
	return &accountRepo{
		db: db,
	}
}

const querySelectBase string = `
		SELECT 
			ta.account_id,
			ta.account_uuid,
			ta.name,
			ta.cpf,
			ta.balance,
			ta.secret,
			ta.created_at
		
		FROM tab_account 				ta
		`

func (r *accountRepo) parseAccount(row scanner) (account entity.Account, err error) {

	err = row.Scan(
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.CPF,
		&account.Balance,
		&account.Secret,
		&account.CreatedAT,
	)

	if err != nil {
		return account, err
	}

	return account, nil
}
func (r *accountRepo) AddTransfer(ctx context.Context, transfer entity.Transfer) (err error) {
	query := `
		INSERT INTO tab_transfer (
			transfer_uuid,
			account_origin_id,
			account_destination_id,
			amount
		) 
		VALUES (?, ?, ?, ?);
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		transfer.TransferUUID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *accountRepo) CreateAccount(ctx context.Context, account entity.Account) (err error) {
	query := `
		INSERT INTO tab_account (
			account_uuid,
			name,
			cpf,
			secret
		) 
		VALUES (?, ?, ?, ?);
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		account.UUID,
		account.Name,
		account.CPF,
		account.Secret,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *accountRepo) GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error) {

	query := querySelectBase + `
		WHERE  	ta.cpf 	= ?
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return account, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(encryptedCPF)

	account, err = r.parseAccount(row)
	if err != nil {
		return account, mysqlutils.HandleMySQLError(err)
	}

	return account, nil
}

func (r *accountRepo) GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error) {

	var params = []interface{}{}

	query := querySelectBase

	var queryCount = `
		SELECT COUNT(*) FROM (
	` + query + `) as count`

	stmt, err := r.db.Prepare(queryCount)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	err = row.Scan(&totalRecords)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return accounts, totalRecords, nil
	}

	if take > 0 {
		query += `
			LIMIT ?
		`
		params = append(params, take)
	}

	if skip > 0 {
		query += `
			OFFSET ?
		`
		params = append(params, skip)
	}

	stmt, err = r.db.Prepare(query)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}
	account := entity.Account{}
	for rows.Next() {
		account, err = r.parseAccount(rows)
		if err != nil {
			return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
		}
		accounts = append(accounts, account)
	}
	return accounts, totalRecords, nil
}

func (r *accountRepo) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {

	var params = []interface{}{}

	query := querySelectBase + `
		WHERE ta.account_uuid = ?
	`
	params = append(params, accountUUID)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return account, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(params...)
	account, err = r.parseAccount(row)
	if err != nil {
		return account, mysqlutils.HandleMySQLError(err)
	}

	return account, nil
}

func (r *accountRepo) GetTransfersByAccountID(ctx context.Context, accountID int64, origin bool) (transfers []entity.Transfer, err error) {
	var params = []interface{}{}

	query := ` 
		SELECT 
			tt.transfer_id,
			tt.transfer_uuid,
			tt.account_origin_id,
			origin.account_uuid,
			tt.account_destination_id,
			dest.account_uuid,
			tt.amount,
			tt.created_at
		
		FROM 	tab_transfer 			tt

		INNER JOIN tab_account origin
			ON origin.account_id = tt.account_origin_id
		
		INNER JOIN tab_account dest
			ON dest.account_id = tt.account_destination_id

	`
	if origin {
		query += `WHERE	tt.account_origin_id 		= 	? `
	} else {
		query += `WHERE	tt.account_destination_id 	= 	? `
	}
	params = append(params, accountID)

	query += `
		ORDER BY tt.created_at desc
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return transfers, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return transfers, mysqlutils.HandleMySQLError(err)
	}

	transfer := entity.Transfer{}
	for rows.Next() {
		err = rows.Scan(
			&transfer.ID,
			&transfer.TransferUUID,
			&transfer.AccountOriginID,
			&transfer.AccountOriginUUID,
			&transfer.AccountDestinationID,
			&transfer.AccountDestinationUUID,
			&transfer.Amount,
			&transfer.CreateAt,
		)
		if err != nil {
			return transfers, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (r *accountRepo) UpdateAccountBalance(ctx context.Context, accountID int64, balance float64) (err error) {

	var params = []interface{}{}
	query := `
		UPDATE 	tab_account
		
		SET 	balance 	= ?

		WHERE  	account_id 	= ?
	`
	params = append(params, balance, accountID)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	return nil
}
