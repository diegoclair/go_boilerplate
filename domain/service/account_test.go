package service

import (
	"context"
	"errors"
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
				mocks.mar.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(result, nil)
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
				mocks.mar.EXPECT().GetAccountByUUID(ctx, args.accountUUID).Times(1).Return(entity.Account{}, errors.New("some error"))
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
