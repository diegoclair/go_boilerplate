package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/account"
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
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(account.Account{
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
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(account.Account{}, errors.New("some error")).Times(1)
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
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(account.Account{
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
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			tt.buildMock(ctx, allMocks, tt.args)

			s := newAuthService(svc)
			_, err := s.Login(ctx, tt.args.cpf, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authService_CreateSession(t *testing.T) {
	type args struct {
		session dto.Session
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		wantErr   bool
	}{
		{
			name: "Should create session without any errors",
			args: args{
				session: dto.Session{
					AccountID:   1,
					SessionUUID: "uuid",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAuthRepo.EXPECT().CreateSession(ctx, args.session).Return(nil).Times(1)
			},
		},
		{
			name: "Should return error when there is some error to create session",
			args: args{
				session: dto.Session{
					AccountID:   1,
					SessionUUID: "uuid",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAuthRepo.EXPECT().CreateSession(ctx, args.session).Return(errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}
			s := newAuthService(svc)
			if err := s.CreateSession(ctx, tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("authService.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authService_GetSessionByUUID(t *testing.T) {
	type args struct {
		sessionUUID string
	}
	tests := []struct {
		name        string
		buildMock   func(ctx context.Context, mocks allMocks, args args)
		args        args
		wantSession dto.Session
		wantErr     bool
	}{
		{
			name: "Should return a session without error",
			args: args{
				sessionUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				result := dto.Session{SessionID: 1, SessionUUID: "123"}
				mocks.mockAuthRepo.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).Return(result, nil).Times(1)
			},
			wantSession: dto.Session{SessionID: 1, SessionUUID: "123"},
			wantErr:     false,
		},
		{
			name: "Should error if database return some error",
			args: args{
				sessionUUID: "123",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAuthRepo.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).Return(dto.Session{}, errors.New("some error")).Times(1)
			},
			wantSession: dto.Session{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			allMocks, svc, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, allMocks, tt.args)
			}
			s := newAuthService(svc)
			gotSession, err := s.GetSessionByUUID(ctx, tt.args.sessionUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.GetSessionByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSession, tt.wantSession) {
				t.Errorf("authService.GetSessionByUUID() = %v, want %v", gotSession, tt.wantSession)
			}
		})
	}
}
