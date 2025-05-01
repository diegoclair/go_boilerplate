package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/crypto"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/util/random"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func createRandomAccount(t *testing.T) entity.Account {
	args := entity.Account{
		UUID: uuid.NewV4().String(),
		Name: random.RandomName(),
		CPF:  random.RandomCPF(),
	}

	var (
		c   = crypto.NewCrypto()
		err error
	)

	args.Password, err = c.HashPassword(random.RandomPassword())
	require.NoError(t, err)

	createID, err := testMysql.Account().CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotZero(t, createID)

	account, err := testMysql.Account().GetAccountByUUID(context.Background(), args.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	validateTwoAccounts(t, args, account)

	return account
}

func validateTwoAccounts(t *testing.T, accountExpected entity.Account, accountToCompare entity.Account) {
	require.Equal(t, float64(0), accountExpected.Balance)
	require.Equal(t, accountExpected.UUID, accountToCompare.UUID)
	require.Equal(t, accountExpected.Name, accountToCompare.Name)
	require.Equal(t, accountExpected.CPF, accountToCompare.CPF)
	require.Equal(t, accountExpected.Password, accountToCompare.Password)
	require.NotZero(t, accountToCompare.ID)
	require.WithinDuration(t, time.Now(), accountToCompare.CreatedAT, 2*time.Second)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestCreateAccountErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newAccountRepo(db).CreateAccount(context.Background(), entity.Account{})
		return err
	})
}

func TestGetAccountByDocument(t *testing.T) {
	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByDocument(context.Background(), account.CPF)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestGetAccountByDocumentErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "account_id", func(db *sql.DB) error {
		_, err := newAccountRepo(db).GetAccountByDocument(context.Background(), "cpf")
		return err
	})
}

func TestGetAccounts(t *testing.T) {
	createRandomAccount(t)
	createRandomAccount(t)

	ctx := context.Background()

	// assert first account created
	accounts, totalRecords, err := testMysql.Account().GetAccounts(ctx, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.LessOrEqual(t, 2, len(accounts))
	require.NotZero(t, totalRecords)

	require.Equal(t, float64(0), accounts[0].Balance)
	require.NotEmpty(t, accounts[0].UUID)
	require.NotEmpty(t, accounts[0].Name)
	require.NotEmpty(t, accounts[0].CPF)
	require.NotEmpty(t, accounts[0].Password)
	require.NotZero(t, accounts[0].ID)
	require.NotZero(t, accounts[0].CreatedAT)

	// assert second account created
	accounts, totalRecords, err = testMysql.Account().GetAccounts(ctx, 10, 1)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.LessOrEqual(t, 1, len(accounts))
	require.NotZero(t, totalRecords)

	require.Equal(t, float64(0), accounts[0].Balance)
	require.NotEmpty(t, accounts[0].UUID)
	require.NotEmpty(t, accounts[0].Name)
	require.NotEmpty(t, accounts[0].CPF)
	require.NotEmpty(t, accounts[0].Password)
	require.NotZero(t, accounts[0].ID)
	require.NotZero(t, accounts[0].CreatedAT)
}

func TestGetAccountsErrorsWithMock(t *testing.T) {
	testForPaginatedSelectErrorsWithMock(t, "account_id", func(db *sql.DB) error {
		_, _, err := newAccountRepo(db).GetAccounts(context.Background(), 10, 0)
		return err
	})
}

func TestGetAccountByUUID(t *testing.T) {
	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByUUID(context.Background(), account.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestGetAccountByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "account_id", func(db *sql.DB) error {
		_, err := newAccountRepo(db).GetAccountByUUID(context.Background(), "uuid")
		return err
	})
}

func TestGetAccountIDByUUID(t *testing.T) {
	account := createRandomAccount(t)

	accountID, err := testMysql.Account().GetAccountIDByUUID(context.Background(), account.UUID)
	require.NoError(t, err)
	require.NotZero(t, accountID)
	require.Equal(t, account.ID, accountID)
}

func TestGetAccountIDByUUIDErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "account_id", func(db *sql.DB) error {
		_, err := newAccountRepo(db).GetAccountIDByUUID(context.Background(), "uuid")
		return err
	})
}

func TestAddTransfer(t *testing.T) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transferUUID := uuid.NewV4().String()
	transferID, err := testMysql.Account().AddTransfer(context.Background(), transferUUID, account.ID, account2.ID, 50)
	require.NoError(t, err)
	require.NotZero(t, transferID)
}

func TestAddTransferErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newAccountRepo(db).AddTransfer(context.Background(), "uuid", 1, 2, 50)
		return err
	})
}

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)

	var balance float64 = 12
	err := testMysql.Account().UpdateAccountBalance(ctx, account.ID, balance)
	require.NoError(t, err)

	updatedAccount, err := testMysql.Account().GetAccountByUUID(ctx, account.UUID)
	require.NoError(t, err)

	require.Equal(t, balance, updatedAccount.Balance)
}

func TestUpdateBalanceErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newAccountRepo(db).UpdateAccountBalance(context.Background(), 1, 12)
	})
}

func TestGetTransfersByAccountID(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	transferUUID := uuid.NewV4().String()
	_, err := testMysql.Account().AddTransfer(context.Background(), transferUUID, account.ID, account2.ID, 50)
	require.NoError(t, err)

	transferUUID = uuid.NewV4().String()

	_, err = testMysql.Account().AddTransfer(context.Background(), transferUUID, account3.ID, account2.ID, 50)
	require.NoError(t, err)

	transfersMade, totalRecordMade, err := testMysql.Account().GetTransfersByAccountID(ctx, account.ID, 0, 0, true)
	require.NoError(t, err)

	require.Len(t, transfersMade, 1)
	require.Equal(t, 1, int(totalRecordMade))

	transfersReceived, totalRecordReceived, err := testMysql.Account().GetTransfersByAccountID(ctx, account2.ID, 10, 0, false)
	require.NoError(t, err)

	require.Len(t, transfersReceived, 2)
	require.Equal(t, 2, int(totalRecordReceived))

	//if skip one record, should return only one record
	transfersReceived, totalRecordReceived, err = testMysql.Account().GetTransfersByAccountID(ctx, account2.ID, 10, 1, false)
	require.NoError(t, err)

	require.Len(t, transfersReceived, 1)
	require.Equal(t, 2, int(totalRecordReceived))
}

func TestGetTransfersByAccountIDErrorsWithMock(t *testing.T) {
	testForPaginatedSelectErrorsWithMock(t, "transfer_id", func(db *sql.DB) error {
		_, _, err := newAccountRepo(db).GetTransfersByAccountID(context.Background(), 1, 0, 0, true)
		return err
	})
}
