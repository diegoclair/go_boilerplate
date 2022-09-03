package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/golang/mock/gomock"
)

func Test_accountService_GetAccountByUUID(t *testing.T) {

	type args struct {
		accountUUID string
	}
	tests := []struct {
		name        string
		buildMock   func(ctx context.Context, mocks repoMock, args args)
		args        args
		wantAccount entity.Account
		wantErr     bool
	}{
		{
			name: "Should return an account without error",
			args: args{
				accountUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
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
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(entity.Account{}, errors.New("some error")).Times(1)
			},
			wantAccount: entity.Account{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			repoMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, repoMocks, tt.args)
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

func Test_accountService_AddBalance(t *testing.T) {

	type args struct {
		accountUUID string
		amount      float64
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks repoMock, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should add balance without any errors",
			args: args{accountUUID: "account123", amount: 7.32},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 50}
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, result.Balance+args.amount).Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should add balance validating floating point",
			args: args{accountUUID: "account123", amount: 0.1},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 0.2}
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, nil).Times(1),
					mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, 0.3).Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error with there is some error to get account by uuid",
			args: args{accountUUID: "account123", amount: 7.32},
			buildMock: func(ctx context.Context, mocks repoMock, args args) {
				result := entity.Account{}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Return(result, fmt.Errorf("some error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			repoMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			s := &accountService{
				svc: svc,
			}

			if tt.buildMock != nil {
				tt.buildMock(ctx, repoMocks, tt.args)
			}
			if err := s.AddBalance(ctx, tt.args.accountUUID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("AddBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
