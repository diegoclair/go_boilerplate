package postgres

import (
	"context"
	"fmt"

	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
)

type accountRepo struct {
	db dbConn
}

func newAccountRepo(db dbConn) contract.AccountRepo {
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
			ta.created_at,
			ta.active

		FROM tab_account 				ta
		`

func (r *accountRepo) parseAccount(row scanner, total ...*int64) (account entity.Account, err error) {
	dests := []any{
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.CPF,
		&account.Balance,
		&account.Password,
		&account.CreatedAT,
		&account.Active,
	}

	if len(total) > 0 && total[0] != nil {
		dests = append(dests, total[0])
	}

	err = row.Scan(dests...)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *accountRepo) AddTransfer(ctx context.Context, transferUUID string, accountOriginID, accountDestinationID int64, amount float64) (transferID int64, err error) {
	query := `
		INSERT INTO tab_transfer (
			transfer_uuid,
			account_origin_id,
			account_destination_id,
			amount
		)
		VALUES ($1, $2, $3, $4)
		RETURNING transfer_id;
	`

	err = r.db.QueryRow(ctx, query,
		transferUUID,
		accountOriginID,
		accountDestinationID,
		amount,
	).Scan(&transferID)
	if err != nil {
		return transferID, handleDBError(err)
	}

	return transferID, nil
}

func (r *accountRepo) CreateAccount(ctx context.Context, account entity.Account) (createdID int64, err error) {
	query := `
		INSERT INTO tab_account (
			account_uuid,
			name,
			cpf,
			secret
		)
		VALUES ($1, $2, $3, $4)
		RETURNING account_id;
	`

	err = r.db.QueryRow(ctx, query,
		account.UUID,
		account.Name,
		account.CPF,
		account.Password,
	).Scan(&createdID)
	if err != nil {
		return createdID, handleDBError(err)
	}

	return createdID, nil
}

func (r *accountRepo) GetAccountByDocument(ctx context.Context, encryptedCPF string) (account entity.Account, err error) {
	query := querySelectBase + `
		WHERE  	ta.cpf 	= $1
	`

	row := r.db.QueryRow(ctx, query, encryptedCPF)
	account, err = r.parseAccount(row)
	if err != nil {
		return account, handleDBError(err)
	}

	return account, nil
}

func (r *accountRepo) GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error) {
	var params = []any{}
	paramIndex := 1

	query := withCount(querySelectBase)

	if take > 0 {
		query += fmt.Sprintf(`
			LIMIT $%d
		`, paramIndex)
		params = append(params, take)
		paramIndex++
	}

	if skip > 0 {
		query += fmt.Sprintf(`
			OFFSET $%d
		`, paramIndex)
		params = append(params, skip)
	}

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return accounts, totalRecords, handleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		account, err := r.parseAccount(rows, &totalRecords)
		if err != nil {
			return accounts, totalRecords, handleDBError(err)
		}

		accounts = append(accounts, account)
	}

	return accounts, totalRecords, nil
}

func (r *accountRepo) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {
	query := querySelectBase + `
		WHERE ta.account_uuid = $1
	`

	row := r.db.QueryRow(ctx, query, accountUUID)
	account, err = r.parseAccount(row)
	if err != nil {
		return account, handleDBError(err)
	}

	return account, nil
}

func (r *accountRepo) GetAccountIDByUUID(ctx context.Context, accountUUID string) (accountID int64, err error) {
	query := `
		SELECT
			account_id

		FROM tab_account
		WHERE account_uuid = $1
	`

	err = r.db.QueryRow(ctx, query, accountUUID).Scan(&accountID)
	if err != nil {
		return accountID, handleDBError(err)
	}

	return accountID, nil
}

func (r *accountRepo) GetTransfersByAccountID(ctx context.Context, accountID, take, skip int64, origin bool) (transfers []entity.Transfer, totalRecords int64, err error) {
	var params = []any{}
	paramIndex := 1

	query := withCount(`
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

	`)

	if origin {
		query += fmt.Sprintf(`WHERE	tt.account_origin_id 		= 	$%d `, paramIndex)
	} else {
		query += fmt.Sprintf(`WHERE	tt.account_destination_id 	= 	$%d `, paramIndex)
	}

	params = append(params, accountID)
	paramIndex++

	query += `
		ORDER BY tt.created_at desc
	`

	if take > 0 {
		query += fmt.Sprintf(`
			LIMIT $%d
		`, paramIndex)
		params = append(params, take)
		paramIndex++
	}

	if skip > 0 {
		query += fmt.Sprintf(`
			OFFSET $%d
		`, paramIndex)
		params = append(params, skip)
	}

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return transfers, totalRecords, handleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var transfer entity.Transfer
		err = rows.Scan(
			&transfer.ID,
			&transfer.TransferUUID,
			&transfer.AccountOriginUUID,
			&transfer.AccountDestinationUUID,
			&transfer.Amount,
			&transfer.CreatedAt,
			&totalRecords,
		)
		if err != nil {
			return transfers, totalRecords, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, totalRecords, nil
}

func (r *accountRepo) UpdateAccountBalance(ctx context.Context, accountID int64, balance float64) (err error) {
	query := `
		UPDATE 	tab_account

		SET 	balance 	= $1,
				update_at 	= NOW()

		WHERE  	account_id 	= $2
	`

	_, err = r.db.Exec(ctx, query, balance, accountID)
	if err != nil {
		return handleDBError(err)
	}

	return nil
}
