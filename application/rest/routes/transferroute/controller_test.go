package transferroute

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

	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

var (
	tokenMaker auth.AuthToken
	onceToken  sync.Once
)

func getTokenMaker(t *testing.T) auth.AuthToken {
	onceToken.Do(func() {
		cfg, err := config.GetConfigEnvironment("../../../../" + config.ConfigDefaultName)
		require.NoError(t, err)

		cfg.App.Auth.AccessTokenDuration = 2 * time.Second
		cfg.App.Auth.RefreshTokenDuration = 2 * time.Second

		tokenMaker, err = auth.NewAuthToken(cfg.App.Auth, logger.NewNoop())
		require.NoError(t, err)
	})
	return tokenMaker
}

func getServerTest(t *testing.T) (transferMock *mocks.MockTransferService, server *echo.Echo, ctrl *gomock.Controller, transferControler *Controller) {
	ctrl = gomock.NewController(t)

	transferMock = mocks.NewMockTransferService(ctrl)
	tokenMaker = getTokenMaker(t)

	transferControler = &Controller{transferMock, routeutils.New(logger.NewNoop())}
	transferRoute := NewRouter(transferControler, "transfers")

	server = echo.New()
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(tokenMaker),
	)

	transferRoute.RegisterRoutes(appGroup, privateGroup)
	return
}

func addAuthorization(ctx context.Context, t *testing.T, req *http.Request, tokenMaker auth.AuthToken, accountUUID, sessionUUID string) {
	token, _, err := tokenMaker.CreateAccessToken(ctx, accountUUID, sessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	req.Header.Set(infra.TokenKey.String(), token)
}

func TestController_handleAddBalance(t *testing.T) {
	ctx := context.Background()

	type args struct {
		body        any
		accountUUID string
		sessionUUID string
	}

	tests := []struct {
		name          string
		args          args
		setupAuth     func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken)
		buildMocks    func(ctx context.Context, transferMock *mocks.MockTransferService, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.TransferReq{
					AccountDestinationUUID: "randomUUID",
					Amount:                 5.55,
				},
				accountUUID: uuid.NewV4().String(),
				sessionUUID: uuid.NewV4().String(),
			},
			setupAuth: func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken) {
				addAuthorization(ctx, t, req, tokenMaker, args.accountUUID, args.sessionUUID)
			},
			buildMocks: func(ctx context.Context, transferMock *mocks.MockTransferService, args args) {
				body := args.body.(viewmodel.TransferReq)
				transferMock.EXPECT().CreateTransfer(ctx,
					entity.Transfer{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).
					Return(nil).MinTimes(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
			},
		},
		{
			name: "Should return error id we have some error on create transfer",
			args: args{
				body: viewmodel.TransferReq{
					AccountDestinationUUID: "randomUUID2",
					Amount:                 8.88,
				},
				accountUUID: uuid.NewV4().String(),
				sessionUUID: uuid.NewV4().String(),
			},
			setupAuth: func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken) {
				addAuthorization(ctx, t, req, tokenMaker, args.accountUUID, args.sessionUUID)
			},
			buildMocks: func(ctx context.Context, mock *mocks.MockTransferService, args args) {
				body := args.body.(viewmodel.TransferReq)
				mock.EXPECT().CreateTransfer(ctx,
					entity.Transfer{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).
					Return(fmt.Errorf("error to create transfer")).MinTimes(1)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusServiceUnavailable, resp.Code)
				require.Contains(t, resp.Body.String(), "Service temporarily unavailable")
				require.Contains(t, resp.Body.String(), "error to create transfer")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			transferMock, server, ctrl, s := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", rootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if tt.setupAuth != nil {
				tt.setupAuth(t, req, tt.args, tokenMaker)
			}

			if tt.buildMocks != nil {
				c := echo.New().NewContext(req, recorder)
				c.Set(infra.AccountUUIDKey.String(), tt.args.accountUUID)
				c.Set(infra.SessionKey.String(), tt.args.sessionUUID)
				ctx := s.utils.Req().GetContext(c)
				tt.buildMocks(ctx, transferMock, tt.args)
			}

			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestController_handleGetTransfers(t *testing.T) {
	ctx := context.Background()

	type args struct {
		accountUUID string
		sessionUUID string
	}
	tests := []struct {
		name          string
		args          args
		buildMocks    func(ctx context.Context, mock *mocks.MockTransferService)
		setupAuth     func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
		sleep         bool
	}{
		{
			name: "Should pass with success",
			args: args{
				accountUUID: uuid.NewV4().String(),
				sessionUUID: uuid.NewV4().String(),
			},
			buildMocks: func(ctx context.Context, mock *mocks.MockTransferService) {
				mock.EXPECT().GetTransfers(ctx).Return([]entity.Transfer{}, nil).Times(1)
			},
			setupAuth: func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken) {
				addAuthorization(ctx, t, req, tokenMaker, args.accountUUID, args.sessionUUID)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resp.Code)
				require.NotEmpty(t, resp.Body)
			},
		},
		{
			name: "Should return expired token",
			setupAuth: func(t *testing.T, req *http.Request, args args, tokenMaker auth.AuthToken) {
				addAuthorization(ctx, t, req, tokenMaker, args.accountUUID, args.sessionUUID)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
				require.Contains(t, resp.Body.String(), "token has expired")
			},
			sleep: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			transferMock, server, ctrl, s := getServerTest(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.setupAuth != nil {
				tt.setupAuth(t, req, tt.args, tokenMaker)
			}

			if tt.buildMocks != nil {
				c := echo.New().NewContext(req, recorder)
				c.Set(infra.AccountUUIDKey.String(), tt.args.accountUUID)
				c.Set(infra.SessionKey.String(), tt.args.sessionUUID)
				ctx := s.utils.Req().GetContext(c)

				tt.buildMocks(ctx, transferMock)
			}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if tt.sleep {
				time.Sleep(2 * time.Second)
			}
			server.ServeHTTP(recorder, req)

			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
