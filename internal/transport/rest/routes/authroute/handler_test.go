package authroute_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/test"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/viewmodel"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_handleLogin(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.Login)

				m.AuthAppMock.EXPECT().Login(ctx, body.ToDto()).Return(entity.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).Return("a123", contract.TokenPayload{}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).Return("r123", contract.TokenPayload{ExpiredAt: time.Now()}, nil).Times(1)
				m.AuthAppMock.EXPECT().CreateSession(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, req dto.Session) error {
						require.NotEmpty(t, req.SessionUUID)
						require.Equal(t, int64(1), req.AccountID)
						require.Equal(t, "r123", req.RefreshToken)
						require.NotEmpty(t, req.RefreshTokenExpiredAt)

						return nil
					},
				).Times(1)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Should return error when body is invalid",
			args: args{
				body: "invalid body",
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error when login fails",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.Login)

				m.AuthAppMock.EXPECT().Login(ctx, body.ToDto()).Return(entity.Account{}, fmt.Errorf("error to login")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to login")
			},
		},
		{
			name: "Should return error when create access token fails",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.Login)

				m.AuthAppMock.EXPECT().Login(ctx, body.ToDto()).Return(entity.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).Return("", contract.TokenPayload{}, fmt.Errorf("error to create access token")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create access token")
			},
		},
		{
			name: "Should return error when create refresh token fails",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.Login)

				m.AuthAppMock.EXPECT().Login(ctx, body.ToDto()).Return(entity.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).Return("a123", contract.TokenPayload{}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).Return("", contract.TokenPayload{}, fmt.Errorf("error to create refresh token")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create refresh token")
			},
		},
		{
			name: "Should return error when create session fails",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "12345678",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				body := args.body.(viewmodel.Login)

				m.AuthAppMock.EXPECT().Login(ctx, body.ToDto()).Return(entity.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, gomock.Any()).Return("a123", contract.TokenPayload{}, nil).Times(1)
				m.AuthTokenMock.EXPECT().CreateRefreshToken(ctx, gomock.Any()).Return("r123", contract.TokenPayload{ExpiredAt: time.Now()}, nil).Times(1)
				m.AuthAppMock.EXPECT().CreateSession(ctx, gomock.Any()).Return(fmt.Errorf("error to create session")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create session")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authroute.Once = sync.Once{}
			authMock, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", authroute.GroupRouteName, authroute.LoginRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, authMock, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleRefreshToken(t *testing.T) {
	type args struct {
		accountUUID string
		sessionUUID string
		body        any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, m test.SvcMocks, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				accountUUID: "aUuid",
				sessionUUID: "sUuid",
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				input := args.body.(viewmodel.RefreshTokenRequest)
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, input.RefreshToken).
					Return(contract.TokenPayload{
						SessionUUID: args.sessionUUID,
						AccountUUID: args.accountUUID,
					}, nil).Times(1)

				ctx = context.WithValue(ctx, infra.AccountUUIDKey, args.accountUUID)
				ctx = context.WithValue(ctx, infra.SessionKey, args.sessionUUID)

				m.AuthAppMock.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).
					Return(dto.Session{
						RefreshTokenExpiredAt: time.Now().Add(2 * time.Hour),
						RefreshToken:          "r123",
					}, nil).Times(1)

				req := contract.TokenPayloadInput{
					AccountUUID: args.accountUUID,
					SessionUUID: args.sessionUUID,
				}
				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, req).
					Return("a123", contract.TokenPayload{}, nil).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "a123")
			},
		},
		{
			name: "Should return error when body is invalid",
			args: args{
				body: "invalid body",
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "Unmarshal type error")
			},
		},
		{
			name: "Should return error when verify token fails",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, gomock.Any()).Return(contract.TokenPayload{}, fmt.Errorf("error to verify token")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to verify token")
			},
		},
		{
			name: "Should return error when get session by uuid fails",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				input := args.body.(viewmodel.RefreshTokenRequest)
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, input.RefreshToken).
					Return(contract.TokenPayload{
						SessionUUID: args.sessionUUID,
						AccountUUID: args.accountUUID,
					}, nil).Times(1)

				ctx = context.WithValue(ctx, infra.AccountUUIDKey, args.accountUUID)
				ctx = context.WithValue(ctx, infra.SessionKey, args.sessionUUID)

				m.AuthAppMock.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).
					Return(dto.Session{}, fmt.Errorf("error to get session by uuid")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to get session by uuid")
			},
		},
		{
			name: "Should return when session token is mismatched",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(contract.TokenPayload{
						SessionUUID: args.sessionUUID,
						AccountUUID: args.accountUUID,
					}, nil).Times(1)

				ctx = context.WithValue(ctx, infra.AccountUUIDKey, args.accountUUID)
				ctx = context.WithValue(ctx, infra.SessionKey, args.sessionUUID)

				m.AuthAppMock.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).
					Return(dto.Session{
						RefreshToken: "r456",
					}, nil).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
				require.Contains(t, resp.Body.String(), "mismatched session token")
			},
		},
		{
			name: "Should return error when refresh token is expired",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(contract.TokenPayload{
						SessionUUID: args.sessionUUID,
						AccountUUID: args.accountUUID,
					}, nil).Times(1)

				ctx = context.WithValue(ctx, infra.AccountUUIDKey, args.accountUUID)
				ctx = context.WithValue(ctx, infra.SessionKey, args.sessionUUID)

				m.AuthAppMock.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).
					Return(dto.Session{
						RefreshTokenExpiredAt: time.Now().Add(-2 * time.Hour),
						RefreshToken:          "r123",
					}, nil).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
				require.Contains(t, resp.Body.String(), "expired session")
			},
		},
		{
			name: "Should return error when create access token fails",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, m test.SvcMocks, args args) {
				m.AuthTokenMock.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(contract.TokenPayload{
						SessionUUID: args.sessionUUID,
						AccountUUID: args.accountUUID,
					}, nil).Times(1)

				ctx = context.WithValue(ctx, infra.AccountUUIDKey, args.accountUUID)
				ctx = context.WithValue(ctx, infra.SessionKey, args.sessionUUID)

				m.AuthAppMock.EXPECT().GetSessionByUUID(ctx, args.sessionUUID).
					Return(dto.Session{
						RefreshTokenExpiredAt: time.Now().Add(2 * time.Hour),
						RefreshToken:          "r123",
					}, nil).Times(1)

				req := contract.TokenPayloadInput{
					AccountUUID: args.accountUUID,
					SessionUUID: args.sessionUUID,
				}

				m.AuthTokenMock.EXPECT().CreateAccessToken(ctx, req).
					Return("", contract.TokenPayload{}, fmt.Errorf("error to create access token")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create access token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", authroute.GroupRouteName, authroute.RefreshTokenRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, false)

			if tt.buildMocks != nil {
				tt.buildMocks(ctx, m, tt.args)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.Echo().ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestHandler_handleLogout(t *testing.T) {
	tests := append(test.PrivateEndpointValidations,
		test.PrivateEndpointTest{
			Name: "Should complete request with no error",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				m.AuthAppMock.EXPECT().Logout(ctx, gomock.Any()).Return(nil).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		test.PrivateEndpointTest{
			Name: "Should return error when logout fails",
			SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m test.SvcMocks) {
				test.AddAuthorization(ctx, t, req, m)
			},
			BuildMocks: func(ctx context.Context, m test.SvcMocks, body any) {
				m.AuthAppMock.EXPECT().Logout(ctx, gomock.Any()).Return(fmt.Errorf("error to logout")).Times(1)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
				require.Contains(t, recorder.Body.String(), "error to logout")
			},
		},
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			authroute.Once = sync.Once{}
			m, server, ctrl := test.GetServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", authroute.GroupRouteName, authroute.LogoutRoute)

			req, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			ctx := test.GetTestContext(t, req, recorder, true)

			if tt.SetupAuth != nil {
				tt.SetupAuth(ctx, t, req, m)
			}

			if tt.BuildMocks != nil {
				tt.BuildMocks(ctx, m, nil)
			}

			server.Echo().ServeHTTP(recorder, req)
			if tt.CheckResponse != nil {
				tt.CheckResponse(t, recorder)
			}
		})
	}
}
