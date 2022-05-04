package mysql

import (
	"testing"
	"time"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/crypto"
	"github.com/diegoclair/go-boilerplate/util/random"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) entity.Account {
	args := entity.Account{
		UUID:   uuid.NewV4().String(),
		Name:   random.RandomName(),
		CPF:    random.RandomCPF(),
		Secret: crypto.GetMd5(random.RandomSecret()),
	}
	err := testMysql.Account().CreateAccount(args)
	require.NoError(t, err)

	account, err := testMysql.Account().GetAccountByUUID(args.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, float64(0), account.Balance)
	require.Equal(t, account.UUID, args.UUID)
	require.Equal(t, account.Name, args.Name)
	require.Equal(t, account.CPF, args.CPF)
	require.Equal(t, account.Secret, args.Secret)
	require.NotZero(t, account.ID)
	require.WithinDuration(t, time.Now(), account.CreatedAT, time.Second)

	return account
}
