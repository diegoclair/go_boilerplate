package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/util/crypto"
	"github.com/diegoclair/go_boilerplate/util/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"go.uber.org/mock/gomock"
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
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the last insert id fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().
					WillReturnResult(sqlmock.NewErrorResult(assert.AnError))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			_, err = repo.CreateAccount(ctx, entity.Account{})
			require.Error(t, err)
		})
	}
}

func TestGetAccountByDocument(t *testing.T) {
	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByDocument(context.Background(), account.CPF)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestGetAccountByDocumentErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"id", "uuid", "name", "cpf", "balance", "password", "created_at"}).
						AddRow(1, "uuid", "name", "cpf", 0, "password", "error"), // invalid string as time.Time
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			_, err = repo.GetAccountByDocument(ctx, "cpf")
			require.Error(t, err)
		})
	}
}

func TestGetAccounts(t *testing.T) {
	createRandomAccount(t)
	createRandomAccount(t)

	// assert first account created
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

	// assert second account created
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

func TestGetAccountsErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "Should not return error when the count returns 0",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(0),
				)
			},
		},
		{
			name: "Should return error when the count prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the count query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"id", "uuid", "name", "cpf", "balance", "password", "created_at"}).AddRow(1, "uuid", "name", "cpf", 0, "password", "2021-01-01 00:00:00"),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			_, _, err = repo.GetAccounts(ctx, 0, 0)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

		})
	}
}

func TestGetAccountByUUID(t *testing.T) {
	account := createRandomAccount(t)

	account2, err := testMysql.Account().GetAccountByDocument(context.Background(), account.CPF)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	validateTwoAccounts(t, account, account2)
}

func TestGetAccountByUUIDErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"id", "uuid", "name", "cpf", "balance", "password", "created_at"}).AddRow(1, "uuid", "name", "cpf", 0, "password", "2021-01-01 00:00:00"),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			_, err = repo.GetAccountByUUID(ctx, "uuid")
			require.Error(t, err)
		})
	}
}

func TestAddTransfer(t *testing.T) {
	account := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transferUUID := uuid.NewV4().String()
	err := testMysql.Account().AddTransfer(context.Background(), transferUUID, account.ID, account2.ID, 50)
	require.NoError(t, err)
}

func TestAddTransferErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			err = repo.AddTransfer(ctx, "uuid", 1, 2, 50)
			require.Error(t, err)
		})
	}
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
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			err = repo.UpdateAccountBalance(ctx, 1, 50)
			require.Error(t, err)
		})
	}
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

func TestGetTransfersByAccountIDErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(dbMock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "Should not return error when the count returns 0",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(0),
				)
			},
		},
		{
			name: "Should return error when the count prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the count query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the prepare fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the query fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the scan fails",
			setupMocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1),
				)

				dbMock.ExpectPrepare("").ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"id", "uuid", "account_origin_id", "account_destination_id", "amount", "created_at"}).AddRow(1, "uuid", 1, 2, 50, "2021-01-01 00:00:00"),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.setupMocks != nil {
				tt.setupMocks(dbMock)
			}

			repo := &accountRepo{
				db: db,
			}

			_, _, err = repo.GetTransfersByAccountID(ctx, 1, 0, 0, true)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
