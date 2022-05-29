package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/crypto"
	"github.com/diegoclair/go-boilerplate/util/random"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func createRandomAccount(t *testing.T) entity.Account {
	args := entity.Account{
		UUID:   uuid.NewV4().String(),
		Name:   random.RandomName(),
		CPF:    random.RandomCPF(),
		Secret: crypto.GetMd5(random.RandomSecret()),
	}
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
	require.Equal(t, accountExpected.Secret, accountToCompare.Secret)
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

	accounts, err := testMysql.Account().GetAccounts(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.LessOrEqual(t, 1, len(accounts))

	require.Equal(t, float64(0), accounts[0].Balance)
	require.NotEmpty(t, accounts[0].UUID)
	require.NotEmpty(t, accounts[0].Name)
	require.NotEmpty(t, accounts[0].CPF)
	require.NotEmpty(t, accounts[0].Secret)
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

	args := entity.Transfer{
		AccountOriginID:      account.ID,
		AccountDestinationID: account2.ID,
		TransferUUID:         uuid.NewV4().String(),
	}
	err := testMysql.Account().AddTransfer(context.Background(), args)
	require.NoError(t, err)
}
