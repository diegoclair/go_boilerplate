package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/entity"
)

func Test_accountService_GetAccountByUUID(t *testing.T) {

	ctx := context.Background()
	repoMocks, svc := newServiceTestMock(t)
	s := &accountService{
		svc: svc,
	}

	type args struct {
		accountUUID string
	}
	tests := []struct {
		name        string
		buildMock   func(ctx context.Context, mocks mocks, args args)
		args        args
		wantAccount entity.Account
		wantErr     bool
	}{
		{
			name: "Should return an account without error",
			args: args{
				accountUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks mocks, args args) {
				result := entity.Account{ID: 1, UUID: "123", Name: "name"}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(result, nil)
			},
			wantAccount: entity.Account{ID: 1, UUID: "123", Name: "name"},
			wantErr:     false,
		},
		{
			name: "Should error if database return some error",
			args: args{
				accountUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks mocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(entity.Account{}, errors.New("some error"))
			},
			wantAccount: entity.Account{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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

	ctx := context.Background()
	repoMocks, svc := newServiceTestMock(t)
	s := &accountService{
		svc: svc,
	}

	type args struct {
		accountUUID string
		amount      float64
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks mocks, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should add balance without any errors",
			args: args{accountUUID: "account123", amount: 7.32},
			buildMock: func(ctx context.Context, mocks mocks, args args) {
				result := entity.Account{ID: 12, UUID: args.accountUUID, Balance: 50}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(result, nil)
				mocks.mockAccountRepo.EXPECT().UpdateAccountBalance(ctx, result.ID, result.Balance+args.amount).Times(1).Return(nil)
			},
		},
		{
			name: "Should return error with there is some error to get account by uuid",
			args: args{accountUUID: "account123", amount: 7.32},
			buildMock: func(ctx context.Context, mocks mocks, args args) {
				result := entity.Account{}
				mocks.mockAccountRepo.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(result, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.buildMock != nil {
				tt.buildMock(ctx, repoMocks, tt.args)
			}
			if err := s.AddBalance(ctx, tt.args.accountUUID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("accountService.AddBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
