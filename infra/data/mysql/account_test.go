package mysql_test

import (
	"context"

	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/util/crypto"
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

	c := crypto.NewCrypto()
	args.Password, _ = c.HashPassword(random.RandomSecret())
	err := testMysql.Account().CreateAccount(context.Background(), args)
	require.NoError(t, err)

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
	require.WithinDuration(t, time.Now(), accountToCompare.CreatedAT, time.Second)
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

	origin := true
	transfersAccount1AsOrigin, err := testMysql.Account().GetTransfersByAccountID(ctx, account.ID, origin)
	require.NoError(t, err)

	require.Len(t, transfersAccount1AsOrigin, 1)

	origin = false
	transfersAccount2AsDestination, err := testMysql.Account().GetTransfersByAccountID(ctx, account2.ID, origin)
	require.NoError(t, err)

	require.Len(t, transfersAccount2AsDestination, 2)
}
