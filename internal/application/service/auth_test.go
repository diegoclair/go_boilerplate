package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"go.uber.org/mock/gomock"
)

func Test_newAuthApp(t *testing.T) {
	m, ctrl := newServiceTestMock(t)
	defer ctrl.Finish()

	want := &authApp{cache: m.mockCacheManager,
		crypto:              m.mockCrypto,
		dm:                  m.mockDataManager,
		log:                 m.mockLogger,
		validator:           m.mockValidator,
		accountSvc:          m.mockAccountSvc,
		accessTokenDuration: time.Minute,
	}

	if got := newAuthApp(m.mockDomain, m.mockAccountSvc, time.Minute); !reflect.DeepEqual(got, want) {
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
				cpf:      "01234567890",
				password: "01234567890",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				gomock.InOrder(
					mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{
						ID:       1,
						UUID:     "uuid",
						Name:     "name",
						CPF:      args.cpf,
						Password: args.password,
						Active:   true,
					}, nil).Times(1),

					mocks.mockCrypto.EXPECT().CheckPassword(args.password, args.password).Return(nil).Times(1),
				)
			},
		},
		{
			name: "Should return error when the account is not active",
			args: args{
				cpf:      "01234567890",
				password: "01234567890",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{
					ID:       1,
					UUID:     "uuid",
					Name:     "name",
					CPF:      args.cpf,
					Password: args.password,
					Active:   false,
				}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error when there is some error to get account by document",
			args: args{
				cpf:      "01234567890",
				password: "01234567890",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).
					Return(entity.Account{}, errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error when the password is wrong",
			args: args{
				cpf:      "01234567890",
				password: "01234567890",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAccountRepo.EXPECT().GetAccountByDocument(ctx, args.cpf).Return(entity.Account{
					ID:       1,
					UUID:     "uuid",
					Name:     "name",
					CPF:      args.cpf,
					Password: args.password,
					Active:   true,
				}, nil).Times(1)

				mocks.mockCrypto.EXPECT().CheckPassword(args.password, args.password).Return(errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Should return error when the input is invalid",
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}

			s := newAuthApp(m.mockDomain, m.mockAccountSvc, time.Minute)

			input := dto.LoginInput{
				CPF:      tt.args.cpf,
				Password: tt.args.password,
			}

			_, err := s.Login(ctx, input)
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
					AccountID:    1,
					SessionUUID:  "d152a340-9a87-4d32-85ad-19df4c9934cd",
					RefreshToken: "token",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAuthRepo.EXPECT().CreateSession(ctx, args.session).Return(int64(0), nil).Times(1)
			},
		},
		{
			name: "Should return error when there is some error to create session",
			args: args{
				session: dto.Session{
					AccountID:    1,
					SessionUUID:  "d152a340-9a87-4d32-85ad-19df4c9934cd",
					RefreshToken: "token",
				},
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockAuthRepo.EXPECT().CreateSession(ctx, args.session).Return(int64(0), errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Should return error when the input is invalid",
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}
			s := newAuthApp(m.mockDomain, m.mockAccountSvc, time.Minute)
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
			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}
			s := newAuthApp(m.mockDomain, m.mockAccountSvc, time.Minute)
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

func Test_authService_Logout(t *testing.T) {
	type args struct {
		accessToken string
	}
	tests := []struct {
		name      string
		buildMock func(ctx context.Context, mocks allMocks, args args)
		args      args
		noSession bool
		wantErr   bool
	}{
		{
			name: "Should logout without any errors",
			args: args{
				accessToken: "token",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockCacheManager.EXPECT().Set(ctx, args.accessToken, "true", gomock.Any()).Return(nil).Times(1)
				mocks.mockAuthRepo.EXPECT().SetSessionAsBlocked(ctx, "session-uuid").Return(nil).Times(1)
			},
		},
		{
			name:      "Should return error when session UUID is not in context",
			args:      args{accessToken: "token"},
			noSession: true,
			wantErr:   true,
		},
		{
			name: "Should return error when there is some error to set string with expiration",
			args: args{
				accessToken: "token",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockCacheManager.EXPECT().Set(ctx, args.accessToken, "true", gomock.Any()).Return(errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Should return error when there is some error to set blocked session",
			args: args{
				accessToken: "token",
			},
			buildMock: func(ctx context.Context, mocks allMocks, args args) {
				mocks.mockCacheManager.EXPECT().Set(ctx, args.accessToken, "true", gomock.Any()).Return(nil).Times(1)
				mocks.mockAuthRepo.EXPECT().SetSessionAsBlocked(ctx, "session-uuid").Return(errors.New("some error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.WithValue(context.Background(), infra.SessionKey, "session-uuid")
			if tt.noSession {
				ctx = context.Background()
			}

			m, ctrl := newServiceTestMock(t)
			defer ctrl.Finish()

			if tt.buildMock != nil {
				tt.buildMock(ctx, m, tt.args)
			}
			s := newAuthApp(m.mockDomain, m.mockAccountSvc, time.Minute)
			if err := s.Logout(ctx, tt.args.accessToken); (err != nil) != tt.wantErr {
				t.Errorf("authService.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
