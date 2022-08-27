package transferroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/mock"
	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

var tokenMaker auth.AuthToken

// var transferMock mocks
//var server *echo.Echo

type mocks struct {
	mapper mapper.Mapper
	mts    *mock.MockTransferService
}

func TestMain(m *testing.M) {
	cfg, err := config.GetConfigEnvironment("../../../../" + config.ConfigDefaultFilepath)
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	cfg.App.Auth.AccessTokenDuration = 2 * time.Second
	cfg.App.Auth.AccessTokenDuration = 2 * time.Second

	tokenMaker, err = auth.NewAuthToken(cfg.App.Auth)
	if err != nil {
		log.Fatal("failed to create authToken: ", err)
	}

	os.Exit(m.Run())
}

func addAuthorization(t *testing.T, req *http.Request, tokenMaker auth.AuthToken) {

	token, _, err := tokenMaker.CreateAccessToken("account123", "session123")
	require.NoError(t, err)

	require.NotEmpty(t, token)

	req.Header.Set(auth.TokenKey.String(), token)
}

func setupTestServer(t *testing.T) (transferMock mocks, server *echo.Echo) {

	ctrl := gomock.NewController(t)
	transferMock = mocks{
		mapper: mapper.New(),
		mts:    mock.NewMockTransferService(ctrl),
	}

	transferControler := NewController(transferMock.mts, transferMock.mapper)
	transferRoute := NewRouter(transferControler, "transfers")

	server = echo.New()
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(tokenMaker),
	)

	transferRoute.RegisterRoutes(appGroup, privateGroup)
	return transferMock, server
}

func TestController_handleAddBalance(t *testing.T) {

	transferMock, server := setupTestServer(t)

	type args struct {
		body any
	}
	tests := []struct {
		name          string
		args          args
		setupAuth     func(t *testing.T, req *http.Request)
		buildMocks    func(ctx context.Context, mock mocks, args args)
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
			buildMocks: func(ctx context.Context, mocks mocks, args args) {
				body := args.body.(viewmodel.TransferReq)
				mocks.mts.EXPECT().CreateTransfer(gomock.Any(), entity.Transfer{AccountDestinationUUID: body.AccountDestinationUUID, Amount: body.Amount}).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Empty(t, resp.Body)
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
