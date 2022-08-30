package transferroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routes/testutil"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

var tokenMaker auth.AuthToken
var transferMock *mock.MockTransferService
var server *echo.Echo

func TestMain(m *testing.M) {
	tokenMaker = testutil.GetTokenMaker()
	transferMock = testutil.NewServiceManagerTest(&testing.T{}).TransferServiceMock

	transferControler := NewController(transferMock, mapper.New())
	transferRoute := NewRouter(transferControler, "transfers")

	server = echo.New()
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(tokenMaker),
	)

	transferRoute.RegisterRoutes(appGroup, privateGroup)

	os.Exit(m.Run())
}

func addAuthorization(t *testing.T, req *http.Request, tokenMaker auth.AuthToken) {

	token, _, err := tokenMaker.CreateAccessToken("account123", "session123")
	require.NoError(t, err)

	require.NotEmpty(t, token)

	req.Header.Set(auth.TokenKey.String(), token)
}

func TestController_handleAddBalance(t *testing.T) {

	type args struct {
		body any
	}
	tests := []struct {
		name          string
		args          args
		setupAuth     func(t *testing.T, req *http.Request)
		buildMocks    func(ctx context.Context, mock *mock.MockTransferService, args args)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Should complete request with no error",
			args: args{
				body: viewmodel.TransferReq{
					AccountDestinationUUID: "randomUUID",
					Amount:                 5.55,
				},
			},
			setupAuth: func(t *testing.T, req *http.Request) {
				addAuthorization(t, req, tokenMaker)
			},
			buildMocks: func(ctx context.Context, mock *mock.MockTransferService, args args) {
				body := args.body.(viewmodel.TransferReq)
				mock.EXPECT().CreateTransfer(gomock.Any(), entity.Transfer{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).Times(1).Return(nil)
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
					AccountDestinationUUID: "randomUUID",
					Amount:                 8.88,
				},
			},
			setupAuth: func(t *testing.T, req *http.Request) {
				addAuthorization(t, req, tokenMaker)
			},
			buildMocks: func(ctx context.Context, mock *mock.MockTransferService, args args) {
				body := args.body.(viewmodel.TransferReq)
				mock.EXPECT().CreateTransfer(gomock.Any(), entity.Transfer{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).Times(1).Return(fmt.Errorf("error to create transfer"))
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

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", rootRoute)

			body, err := json.Marshal(tt.args.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			if tt.setupAuth != nil {
				tt.setupAuth(t, req)
			}

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
				tt.buildMocks(ctx, transferMock, tt.args)
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.ServeHTTP(recorder, req)
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestController_handleGetTransfers(t *testing.T) {

	tests := []struct {
		name          string
		buildMocks    func(ctx context.Context, mock *mock.MockTransferService)
		setupAuth     func(t *testing.T, req *http.Request)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
		sleep         bool
	}{
		{
			name: "Should pass with success",
			buildMocks: func(ctx context.Context, mock *mock.MockTransferService) {
				mock.EXPECT().GetTransfers(gomock.Any()).Times(1).Return([]entity.Transfer{}, nil)
			},
			setupAuth: func(t *testing.T, req *http.Request) {
				addAuthorization(t, req, tokenMaker)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, resp.Code)
				require.NotEmpty(t, resp.Body)
			},
		},

		{
			name: "Should return expired token",
			buildMocks: func(ctx context.Context, mock *mock.MockTransferService) {
				mock.EXPECT().GetTransfers(gomock.Any()).Times(1).Return([]entity.Transfer{}, nil)
			},
			setupAuth: func(t *testing.T, req *http.Request) {
				addAuthorization(t, req, tokenMaker)
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

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transfers%s", rootRoute)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tt.setupAuth != nil {
				tt.setupAuth(t, req)
			}

			if tt.buildMocks != nil {
				e := echo.New()
				ctx := routeutils.GetContext(e.NewContext(req, recorder))
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
