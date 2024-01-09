package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/golang/mock/gomock"
)

func Test_newAuthService(t *testing.T) {
	_, svc, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()

	want := &authService{svc: svc}

	if got := newAuthService(svc); !reflect.DeepEqual(got, want) {
		t.Errorf("newAuthService() = %v, want %v", got, want)
	}
}

func Test_authService_Login(t *testing.T) {
	type args struct {
		cpf      string
		password string
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should login without any errors",
			args: args{
				cpf:      "123",
				password: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{
						ID:       1,
						UUID:     "uuid",
						Name:     "name",
						CPF:      args.cpf,
						Password: "123",
					}, nil).Times(1),

					mocks.mockCrypto.EXPECT().CheckPassword(args.password, "123").Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error when there is some error to get account by document",
			args: args{
				cpf:      "123",
				password: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{}, errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the password is wrong",
			args: args{
				cpf:      "123",
				password: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{
					ID:       1,
					UUID:     "uuid",
					Name:     "name",
					CPF:      args.cpf,
					Password: "123",
				}, nil).Times(1)

				mocks.mockCrypto.EXPECT().CheckPassword(args.password, "123").Return(errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			allMockss, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			tt.buildMock(ctx, allMockss, tt.args)

			s := newAuthService(svc)
			_, err := s.Login(ctx, tt.args.cpf, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
