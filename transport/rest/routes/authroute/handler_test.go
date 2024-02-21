package authroute

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

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/goswag"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mock struct {
	authService *mocks.MockAuthService
	authToken   *mocks.MockAuthToken
}

func getServerTest(t *testing.T) (authMock mock, server goswag.Echo, ctrl *gomock.Controller, accountHandler *Handler) {
	ctrl = gomock.NewController(t)
	authMock = mock{
		authService: mocks.NewMockAuthService(ctrl),
		authToken:   mocks.NewMockAuthToken(ctrl),
	}

	v, err := validator.NewValidator()
	require.NoError(t, err)

	accountHandler = NewHandler(authMock.authService, authMock.authToken, routeutils.New(), v)
	accountRoute := NewRouter(accountHandler, RouteName)

	server = goswag.NewEcho()
	appGroup := server.Group("/")
	g := &routeutils.EchoGroups{
		AppGroup: appGroup,
	}

	accountRoute.RegisterRoutes(g)
	return
}

func TestHandler_handleLogin(t *testing.T) {
	type args struct {
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock mock, args args)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				mock.authToken.EXPECT().CreateAccessToken(ctx, "uuid", gomock.Any()).Return("a123", &auth.TokenPayload{}, nil).Times(1)
				mock.authToken.EXPECT().CreateRefreshToken(ctx, "uuid", gomock.Any()).Return("r123", &auth.TokenPayload{ExpiredAt: time.Now()}, nil).Times(1)
				mock.authService.EXPECT().CreateSession(ctx, gomock.Any()).DoAndReturn(
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
			name: "Should return error when cpf is invalid",
			args: args{
				body: viewmodel.Login{
					CPF:      "0123456789",
					Password: "12345678",
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "The field 'CPF' should be a valid cpf")
			},
		},
		{
			name: "Should return error when password is invalid",
			args: args{
				body: viewmodel.Login{
					CPF:      "01234567890",
					Password: "1234567",
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "The field 'Password' should have the minimun lenhgt or value: 8")
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{}, fmt.Errorf("error to login")).Times(1)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				mock.authToken.EXPECT().CreateAccessToken(ctx, "uuid", gomock.Any()).Return("", nil, fmt.Errorf("error to create access token")).Times(1)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				mock.authToken.EXPECT().CreateAccessToken(ctx, "uuid", gomock.Any()).Return("a123", &auth.TokenPayload{}, nil).Times(1)
				mock.authToken.EXPECT().CreateRefreshToken(ctx, "uuid", gomock.Any()).Return("", nil, fmt.Errorf("error to create refresh token")).Times(1)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.Login)
				mock.authService.EXPECT().Login(ctx, input.CPF, input.Password).Return(account.Account{ID: 1, UUID: "uuid"}, nil).Times(1)
				mock.authToken.EXPECT().CreateAccessToken(ctx, "uuid", gomock.Any()).Return("a123", &auth.TokenPayload{}, nil).Times(1)
				mock.authToken.EXPECT().CreateRefreshToken(ctx, "uuid", gomock.Any()).Return("r123", &auth.TokenPayload{ExpiredAt: time.Now()}, nil).Times(1)
				mock.authService.EXPECT().CreateSession(ctx, gomock.Any()).Return(fmt.Errorf("error to create session")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create session")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			authMock, server, ctrl, authHandler := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s%s", RouteName, loginRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := authHandler.utils.Req().GetContext(e.NewContext(req, recorder))
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
		body any
	}

	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock mock, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				input := args.body.(viewmodel.RefreshTokenRequest)
				mock.authToken.EXPECT().VerifyToken(ctx, input.RefreshToken).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
					Return(dto.Session{
						RefreshTokenExpiredAt: time.Now().Add(2 * time.Hour),
						RefreshToken:          "r123",
					}, nil).Times(1)

				mock.authToken.EXPECT().CreateAccessToken(ctx, "aUuid", "sUuid").
					Return("a123", &auth.TokenPayload{}, nil).Times(1)
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
			name: "Should return error when refresh token is invalid",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "",
				},
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, resp.Code)
				require.Contains(t, resp.Body.String(), "The field 'RefreshToken' is required")
			},
		},
		{
			name: "Should return error when verify token fails",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).Return(nil, fmt.Errorf("error to verify token")).Times(1)
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
					Return(dto.Session{}, fmt.Errorf("error to get session by uuid")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to get session by uuid")
			},
		},
		{
			name: "Should return when session is blocked",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
					Return(dto.Session{
						IsBlocked: true,
					}, nil).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
				require.Contains(t, resp.Body.String(), "blocked session")
			},
		},
		{
			name: "Should return when session token is mismatched",
			args: args{
				body: viewmodel.RefreshTokenRequest{
					RefreshToken: "r123",
				},
			},
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
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
			buildMocks: func(ctx context.Context, mock mock, args args) {
				mock.authToken.EXPECT().VerifyToken(ctx, gomock.Any()).
					Return(&auth.TokenPayload{
						SessionUUID: "sUuid",
						AccountUUID: "aUuid",
					}, nil).Times(1)

				mock.authService.EXPECT().GetSessionByUUID(ctx, "sUuid").
					Return(dto.Session{
						RefreshTokenExpiredAt: time.Now().Add(2 * time.Hour),
						RefreshToken:          "r123",
					}, nil).Times(1)

				mock.authToken.EXPECT().CreateAccessToken(ctx, "aUuid", "sUuid").
					Return("", nil, fmt.Errorf("error to create access token")).Times(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "error to create access token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			once = sync.Once{}
			authMock, server, ctrl, authHandler := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/refresh-token", RouteName)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := authHandler.utils.Req().GetContext(e.NewContext(req, recorder))
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
