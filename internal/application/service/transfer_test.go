package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/util/number"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_newTransferService(t *testing.T) {
	m, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()

	want := &transferService{dm: m.mockDataManager, accountSvc: m.mockAccountSvc, log: m.mockLogger, validator: m.mockValidator}

	if got := newTransferService(m.mockDomain, m.mockAccountSvc); !reflect.DeepEqual(got, want) {
		t.Errorf("newTransferService() = %v, want %v", got, want)
	}
}

func Test_transferService_CreateTransfer(t *testing.T) {
	type args struct {
		accountUUIDFromContext string
		transfer               dto.TransferInput
	}

	tests := []struct {
		name      string
		args      args
		buildMock func(ctx context.Context, mocks allMocks, args args)
		wantErr   bool
	}{
		{
			name: "Should pass without error",
			args: args{
				accountUUIDFromContext: "account-from-123",
				transfer: dto.TransferInput{
					Amount:                 5.0,
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{
							ID:      1,
							Balance: 10.50,
						}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{
							ID:      2,
							Balance: 25.50,
						}, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(r contract.DataManager) error) error {
							return fn(mocks.mockDataManager)
						},
					).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(int64(0), nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), 5.50).
						Return(nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(2), 30.50).
						Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error if accountUUID does not exists on database",
			args: args{
				accountUUIDFromContext: "account-non-exists",
				transfer: dto.TransferInput{
					Amount:                 5.0,
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
					Return(entity.Account{}, fmt.Errorf("account not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error if from account has not sufficient balance to transfer",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					Amount:                 18,
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
					Return(entity.Account{
						Balance: 15,
					}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should not receive value with floating point problem",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					Amount:                 0.1,
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(

					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{
							ID:      1,
							Balance: 0.3,
						}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{
							ID:      2,
							Balance: 0.2,
						}, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(r contract.DataManager) error) error {
							return fn(mocks.mockDataManager)
						},
					).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(int64(0), nil).Times(1),
					//if we remove the number.RoundFloat of destination balance, here we would have 0.19999999999999998 instead of 0.2
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), 0.2).
						Return(nil).Times(1),
					//if we remove the number.RoundFloat of destination balance, here we would have 0.30000000000000004 instead of 0.3
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(2), 0.3).
						Return(nil).Times(1),
				)
			},
			wantErr: false,
		},
		{
			name: "Should return error if there is some error to get account by uuid",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
					Return(entity.Account{}, assert.AnError).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to get account destination by uuid",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{Balance: 5}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{}, assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if the destination account is not found",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{ID: 1, Balance: 5}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{}, fmt.Errorf("no rows in result set")).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if the destination account is the same as the origin account",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{ID: 1, Balance: 5}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{ID: 1}, nil).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to begin transaction",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{ID: 1, Balance: 5}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{ID: 2}, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).Return(assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to add transfer",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(entity.Account{ID: 1, Balance: 5}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{ID: 2}, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(r contract.DataManager) error) error {
							return fn(mocks.mockDataManager)
						},
					).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(int64(0), assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to update origin account balance",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				respAccountFrom := entity.Account{ID: 1, Balance: 4}
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(respAccountFrom, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{ID: 2}, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(r contract.DataManager) error) error {
							return fn(mocks.mockDataManager)
						},
					).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(int64(0), nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), number.RoundFloat(respAccountFrom.Balance-args.transfer.Amount, 2)).
						Return(assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to update destination account balance",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: dto.TransferInput{
					AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					Amount:                 2,
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				respAccountFrom := entity.Account{ID: 1, Balance: 4}
				respAccountDest := entity.Account{ID: 2, Balance: 5}
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccount(ctx).
						Return(respAccountFrom, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(respAccountDest, nil).Times(1),
					mocks.mockDataManager.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(r contract.DataManager) error) error {
							return fn(mocks.mockDataManager)
						},
					).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(int64(0), nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), number.RoundFloat(respAccountFrom.Balance-args.transfer.Amount, 2)).
						Return(nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(2), number.RoundFloat(respAccountDest.Balance+args.transfer.Amount, 2)).
						Return(assert.AnError).Times(1),
				)
			},
			wantErr: true,
		},
		{
			name:    "Should return error if we have invalid transfer input",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := newTransferService(m.mockDomain, m.mockAccountSvc)

			if tt.args.accountUUIDFromContext != "" {
				ctx = context.WithValue(ctx, infra.AccountUUIDKey, tt.args.accountUUIDFromContext)
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}

			if err := s.CreateTransfer(ctx, tt.args.transfer); (err != nil) != tt.wantErr {
				t.Errorf("transferService.CreateTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transferService_GetTransfers(t *testing.T) {
	type args struct {
		accountUUIDFromContext string
	}
	tests := []struct {
		name      string
		args      args
		buildMock func(ctx context.Context, mocks allMocks, args args)
		wantErr   bool
	}{
		{
			name: "Should pass without error",
			args: args{
				accountUUIDFromContext: "account-123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountSvc.EXPECT().GetLoggedAccountID(ctx).
						Return(int64(1), nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetTransfersByAccountID(ctx, int64(1), int64(0), int64(0), true).
						Return([]entity.Transfer{}, int64(0), nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetTransfersByAccountID(ctx, int64(1), int64(0), int64(0), false).
						Return([]entity.Transfer{}, int64(0), nil).Times(1),
				)
			},
		},
		{
			name: "Should return error if there is some error to get account by uuid",
			args: args{
				accountUUIDFromContext: "account-123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccountID(ctx).
					Return(int64(0), assert.AnError).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to get transfers made by account id",
			args: args{
				accountUUIDFromContext: "account-123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccountID(ctx).
					Return(int64(1), nil).Times(1)
				mocks.mockAccountRepo.EXPECT().GetTransfersByAccountID(ctx, int64(1), int64(0), int64(0), true).
					Return([]entity.Transfer{}, int64(0), assert.AnError).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error if there is some error to get transfers received by account id",
			args: args{
				accountUUIDFromContext: "account-123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountSvc.EXPECT().GetLoggedAccountID(ctx).
					Return(int64(1), nil).Times(1)
				mocks.mockAccountRepo.EXPECT().GetTransfersByAccountID(ctx, int64(1), int64(0), int64(0), true).
					Return([]entity.Transfer{}, int64(0), nil).Times(1)
				mocks.mockAccountRepo.EXPECT().GetTransfersByAccountID(ctx, int64(1), int64(0), int64(0), false).
					Return([]entity.Transfer{}, int64(0), assert.AnError).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := newTransferService(m.mockDomain, m.mockAccountSvc)

			if tt.args.accountUUIDFromContext != "" {
				ctx = context.WithValue(ctx, infra.AccountUUIDKey, tt.args.accountUUIDFromContext)
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}

			if _, _, err := s.GetTransfers(ctx, 0, 0); (err != nil) != tt.wantErr {
				t.Errorf("transferService.GetTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
