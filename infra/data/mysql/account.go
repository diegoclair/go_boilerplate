package mysql

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/domain/transfer"
	"github.com/diegoclair/go_utils/mysqlutils"
)

type accountRepo struct {
	db dbConnection
}

func newAccountRepo(db dbConnection) contract.AccountRepo {
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

func (r *accountRepo) parseAccount(row scanner) (account account.Account, err error) {

	err = row.Scan(
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.CPF,
		&account.Balance,
		&account.Password,
		&account.CreatedAT,
	)

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *accountRepo) AddTransfer(ctx context.Context, transferUUID string, accountOriginID, accountDestinationID int64, amount float64) (err error) {
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
		transferUUID,
		accountOriginID,
		accountDestinationID,
		amount,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *accountRepo) CreateAccount(ctx context.Context, account account.Account) (createdID int64, err error) {
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
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		account.UUID,
		account.Name,
		account.CPF,
		account.Password,
	)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	createdID, err = result.LastInsertId()
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	return createdID, nil
}

func (r *accountRepo) GetAccountByDocument(ctx context.Context, encryptedCPF string) (account account.Account, err error) {

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

func (r *accountRepo) GetAccounts(ctx context.Context, take, skip int64) (accounts []account.Account, totalRecords int64, err error) {

	var params = []interface{}{}

	query := querySelectBase

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
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

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	var account account.Account
	for rows.Next() {
		account, err = r.parseAccount(rows)
		if err != nil {
			return accounts, totalRecords, mysqlutils.HandleMySQLError(err)
		}

		accounts = append(accounts, account)
	}

	return accounts, totalRecords, nil
}

func (r *accountRepo) GetAccountByUUID(ctx context.Context, accountUUID string) (account account.Account, err error) {

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

func (r *accountRepo) GetTransfersByAccountID(ctx context.Context, accountID, take, skip int64, origin bool) (transfers []transfer.Transfer, totalRecords int64, err error) {
	var params = []interface{}{}

	query := ` 
		SELECT 
			tt.transfer_id,
			tt.transfer_uuid,
			origin.account_uuid AS account_origin_uuid,
			dest.account_uuid AS account_destination_uuid,
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

	totalRecords, err = getTotalRecordsFromQuery(ctx, r.db, query, params...)
	if err != nil {
		return transfers, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	if totalRecords < 1 {
		return transfers, totalRecords, nil
	}

	query += `
		ORDER BY tt.created_at desc
	`

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

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return transfers, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return transfers, totalRecords, mysqlutils.HandleMySQLError(err)
	}

	transfer := transfer.Transfer{}
	for rows.Next() {
		err = rows.Scan(
			&transfer.ID,
			&transfer.TransferUUID,
			&transfer.AccountOriginUUID,
			&transfer.AccountDestinationUUID,
			&transfer.Amount,
			&transfer.CreatedAt,
		)
		if err != nil {
			return transfers, totalRecords, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, totalRecords, nil
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
