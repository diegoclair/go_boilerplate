package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_newAccountService(t *testing.T) {
	_, svc, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()

	want := &accountService{svc: svc}

	if got := newAccountService(svc); !reflect.DeepEqual(got, want) {
		t.Errorf("newAccountService() = %v, want %v", got, want)
	}
}

func Test_accountService_CreateAccount(t *testing.T) {
	type args struct {
		account dto.AccountInput
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should create account without any errors",
			args: args{account: dto.AccountInput{
				Name:     "name",
				CPF:      "01234567890",
				Password: "01234567890",
			}},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.account.CPF).Return(entity.Account{}, errors.New("No records find")).Times(1),
					mocks.mockCrypto.EXPECT().HashPassword(args.account.Password).Return("123", nil).Times(1),

					mocks.mockAccountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, account entity.Account) (int64, error) {
						require.Equal(t, args.account.Name, account.Name)
						require.Equal(t, args.account.CPF, account.CPF)
						require.NotEmpty(t, account.Password)
						return int64(0), nil
					}).Times(1),
				)
			},
		},
		{
			name: "Should return error with there is some error to get account by document",
			args: args{account: dto.AccountInput{
				Name:     "name",
				CPF:      "01234567890",
				Password: "01234567890",
			}},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.account.CPF).Return(entity.Account{}, errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error with there is some error to create account",
			args: args{account: dto.AccountInput{
				Name:     "name",
				CPF:      "01234567890",
				Password: "01234567890",
			}},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.account.CPF).Return(entity.Account{}, errors.New("No records find")).Times(1),
					mocks.mockCrypto.EXPECT().HashPassword(args.account.Password).Return("123", nil).Times(1),
					mocks.mockAccountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).Return(int64(0), errors.New("some error")).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error with there is some error to hash password",
			args: args{account: dto.AccountInput{
				Name:     "name",
				CPF:      "01234567890",
				Password: "01234567890",
			}},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.account.CPF).Return(entity.Account{}, errors.New("No records find")).Times(1),
					mocks.mockCrypto.EXPECT().HashPassword(args.account.Password).Return("", errors.New("some error")).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error cpf already in use",
			args: args{account: dto.AccountInput{
				Name:     "name",
				CPF:      "01234567890",
				Password: "01234567890",
			}},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.account.CPF).Return(entity.Account{}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Should return error with invalid input",
			args:    args{account: dto.AccountInput{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}
			if err := s.CreateAccount(ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountService_AddBalance(t *testing.T) {
	type args struct {
		accountUUID string
		amount      float64
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should add balance without any errors",
			args: args{accountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd", amount: 7.32},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 50}
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, result.Balance+args.amount).Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should add balance validating floating point",
			args: args{accountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd", amount: 0.1},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 0.2}
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, 0.3).Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error with there is some error to get account by uuid",
			args: args{accountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd", amount: 7.32},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := entity.Account{}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, assert.AnError).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error with there is some error to update account balance",
			args: args{accountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd", amount: 7.32},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 50}
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, result.Balance+args.amount).Return(assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name:    "Should return error with invalid input",
			args:    args{accountUUID: "", amount: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}

			input := dto.AddBalanceInput{
				AccountUUID: tt.args.accountUUID,
				Amount:      tt.args.amount,
			}

			if err := s.AddBalance(ctx, input); (err != nil) != tt.wantErr {
				t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountService_GetAccounts(t *testing.T) {
	type args struct {
		take int64
		skip int64
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		want      []entity.Account
		want1     int64
		wantErr   bool
	}{
		{
			name: "Should return accounts without any errors",
			args: args{take: 10, skip: 0},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := []entity.Account{{ID: 1, UUID: "123", Name: "name"}}
				mocks.mockAccountRepo.EXPECT().GetAccounts(ctx, args.take, args.skip).Return(result, int64(1), nil).Times(1)
			},
			want:    []entity.Account{{ID: 1, UUID: "123", Name: "name"}},
			want1:   1,
			wantErr: false,
		},
		{
			name: "Should return error with there is some error to get accounts",
			args: args{take: 10, skip: 0},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccounts(ctx, args.take, args.skip).Return([]entity.Account{}, int64(0), errors.New("some error")).Times(1)
			},
			want:    []entity.Account{},
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}

			got, got1, err := s.GetAccounts(ctx, tt.args.take, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountService.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accountService.GetAccounts() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("accountService.GetAccounts() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_accountService_GetAccountByUUID(t *testing.T) {

	type args struct {
		accountUUID string
	}
	tests := []struct {
		name        string
		buildMock   func(ctx context.Context, mocks allMocks, args args)
		args        args
		wantAccount entity.Account
		wantErr     bool
	}{
		{
			name: "Should return an account without error",
			args: args{
				accountUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := entity.Account{ID: 1, UUID: "123", Name: "name"}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1)
			},
			wantAccount: entity.Account{ID: 1, UUID: "123", Name: "name"},
			wantErr:     false,
		},
		{
			name: "Should error if database return some error",
			args: args{
				accountUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(entity.Account{}, errors.New("some error")).Times(1)
			},
			wantAccount: entity.Account{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}

			gotAccount, err := s.GetAccountByUUID(ctx, tt.args.accountUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountService.GetAccountByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAccount, tt.wantAccount) {
				t.Errorf("accountService.GetAccountByUUID() = %v, want %v", gotAccount, tt.wantAccount)
			}
		})
	}
}
