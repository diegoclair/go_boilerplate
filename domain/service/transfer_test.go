package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/golang/mock/gomock"
)

func Test_transferService_CreateTransfer(t *testing.T) {

	type args struct {
		accountUUIDFromContext string
		transfer               entity.Transfer
	}

	tests := []struct {
		name      string
		args      args
		buildMock func(ctx context.Context, mocks repoMock, args args)
		wantErr   bool
	}{
		{
			name: "Should pass without error",
			args: args{
				accountUUIDFromContext: "account-from-123",
				transfer: entity.Transfer{
					Amount:                 5.0,
					AccountDestinationUUID: "account-dest-123",
				},
			},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUIDFromContext).
						Return(entity.Account{
							ID:      1,
							Balance: 10.50,
						}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{
							ID:      2,
							Balance: 25.50,
						}, nil).Times(1),
					mocks.mockDataManager.EXPECT().Begin().Return(mocks.mockTransaction, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), 5.50).
						Return(nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(2), 30.50).
						Return(nil).Times(1),
					mocks.mockTransaction.EXPECT().Commit().Return(nil).Times(1),
					mocks.mockTransaction.EXPECT().Rollback().Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error if accountUUID from context is empty",
			args: args{
				transfer: entity.Transfer{},
			},
			wantErr: true,
		},
		{
			name: "Should return error if accountUUID does not exists on database",
			args: args{
				accountUUIDFromContext: "account-non-exists",
				transfer:               entity.Transfer{},
			},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUIDFromContext).
					Return(entity.Account{}, fmt.Errorf("account not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error if from account has not sufficient balance to transfer",
			args: args{
				accountUUIDFromContext: "account-123",
				transfer: entity.Transfer{
					Amount: 18,
				},
			},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUIDFromContext).
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
				transfer: entity.Transfer{
					Amount: 0.1,
				},
			},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				gomock.InOrder(

					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUIDFromContext).
						Return(entity.Account{
							ID:      1,
							Balance: 0.3,
						}, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.transfer.AccountDestinationUUID).
						Return(entity.Account{
							ID:      2,
							Balance: 0.2,
						}, nil).Times(1),
					mocks.mockDataManager.EXPECT().Begin().Return(mocks.mockTransaction, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().AddTransfer(ctx, gomock.Not(""), int64(1), int64(2), args.transfer.Amount).
						Return(nil).Times(1),
					//if we remove the number.RoundFloat of destination balance, here we would have 0.19999999999999998 instead of 0.2
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(1), 0.2).
						Return(nil).Times(1),
					//if we remove the number.RoundFloat of destination balance, here we would have 0.30000000000000004 instead of 0.3
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, int64(2), 0.3).
						Return(nil).Times(1),
					mocks.mockTransaction.EXPECT().Commit().Return(nil).Times(1),
					mocks.mockTransaction.EXPECT().Rollback().Return(nil).Times(1),
				)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			repoMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &transferService{
				svc: svc,
			}

			if tt.args.accountUUIDFromContext != "" {
				ctx = context.WithValue(ctx, auth.AccountUUIDKey, tt.args.accountUUIDFromContext)
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, repoMocks, tt.args)
			}

			if err := s.CreateTransfer(ctx, tt.args.transfer); (err != nil) != tt.wantErr {
				t.Errorf("transferService.CreateTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
