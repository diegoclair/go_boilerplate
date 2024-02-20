package mysql_test

import (
	"context"

	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/util/crypto"
	"github.com/diegoclair/go_boilerplate/util/random"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func createRandomAccount(t *testing.T) account.Account {
	args := account.Account{
		UUID: uuid.NewV4().String(),
		Name: random.RandomName(),
		CPF:  random.RandomCPF(),
	}

	var (
		c   = crypto.NewCrypto()
		err error
	)

	args.Password, err = c.HashPassword(random.RandomSecret())
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

func validateTwoAccounts(t *testing.T, accountExpected account.Account, accountToCompare account.Account) {
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

func TestGetAccountByDocument(t *testing.T) {
	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByDocument(context.Background(), account.CPF)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestGetAccounts(t *testing.T) {
	createRandomAccount(t)
	createRandomAccount(t)

	accounts, totalRecords, err := testMysql.Account().GetAccounts(context.Background(), 10, 0)
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

	accounts, totalRecords, err = testMysql.Account().GetAccounts(context.Background(), 10, 1)
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

func TestGetAccountByUUID(t *testing.T) {

	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByDocument(context.Background(), account.CPF)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestAddTransfer(t *testing.T) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transferUUID := uuid.NewV4().String()
	err := testMysql.Account().AddTransfer(context.Background(), transferUUID, account.ID, account2.ID, 50)
	require.NoError(t, err)
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

func TestGetTransfersByAccountID(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	transferUUID := uuid.NewV4().String()
	err := testMysql.Account().AddTransfer(context.Background(), transferUUID, account.ID, account2.ID, 50)
	require.NoError(t, err)

	transferUUID = uuid.NewV4().String()

	err = testMysql.Account().AddTransfer(context.Background(), transferUUID, account3.ID, account2.ID, 50)
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
