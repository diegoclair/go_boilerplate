package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	infraMocks "github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/internal/transport/rest/serverMiddleware"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/goswag"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

type SvcMocks struct {
	AccountAppMock  *mocks.MockAccountApp
	AuthAppMock     *mocks.MockAuthApp
	AuthTokenMock   *infraMocks.MockAuthToken
	CacheMock       *mocks.MockCacheManager
	TransferAppMock *mocks.MockTransferApp
}

func GetServerTest(t *testing.T) (m SvcMocks, server goswag.Echo, ctrl *gomock.Controller) {
	t.Helper()

	ctrl = gomock.NewController(t)
	m = SvcMocks{
		AccountAppMock:  mocks.NewMockAccountApp(ctrl),
		AuthAppMock:     mocks.NewMockAuthApp(ctrl),
		AuthTokenMock:   infraMocks.NewMockAuthToken(ctrl),
		CacheMock:       mocks.NewMockCacheManager(ctrl),
		TransferAppMock: mocks.NewMockTransferApp(ctrl),
	}

	server = goswag.NewEcho()
	server.Echo().HTTPErrorHandler = func(err error, c echo.Context) {
		_ = routeutils.HandleError(c, err)
	}
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(getTestTokenMaker(t), m.CacheMock),
	)

	g := &routeutils.EchoGroups{
		AppGroup:     appGroup,
		PrivateGroup: privateGroup,
	}

	accountHandler := accountroute.NewHandler(m.AccountAppMock)
	accountRoute := accountroute.NewRouter(accountHandler)
	authHandler := authroute.NewHandler(m.AuthAppMock, m.AuthTokenMock)
	authRoute := authroute.NewRouter(authHandler)
	transferHandler := transferroute.NewHandler(m.TransferAppMock)
	transferRoute := transferroute.NewRouter(transferHandler)

	accountRoute.RegisterRoutes(g)
	authRoute.RegisterRoutes(g)
	transferRoute.RegisterRoutes(g)
	return
}

var (
	tokenMaker contract.AuthToken
	onceToken  sync.Once
)

func getTestTokenMaker(t *testing.T) contract.AuthToken {
	t.Helper()

	onceToken.Do(func() {
		cfg := configmock.New()
		var err error

		cfg.Auth.AccessTokenDuration = 2 * time.Second
		cfg.Auth.RefreshTokenDuration = 2 * time.Second

		tokenMaker, err = auth.NewAuthToken(cfg.Auth.AccessTokenDuration,
			cfg.Auth.RefreshTokenDuration,
			cfg.Auth.PasetoSymmetricKey,
			cfg.GetLogger(),
		)
		require.NoError(t, err)
		require.NotEmpty(t, tokenMaker)
	})
	return tokenMaker
}

var (
	accountUUID = uuid.Must(uuid.NewV7()).String()
	sessionUUID = uuid.Must(uuid.NewV7()).String()
)

func AddAuthorization(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
	t.Helper()

	token := addAuthorizationWithNoCache(ctx, t, req)
	m.CacheMock.EXPECT().GetString(gomock.Any(), token).Return("", nil).Times(1)
}

func addAuthorizationWithNoCache(ctx context.Context, t *testing.T, req *http.Request) (token string) {
	t.Helper()

	tokenMaker := getTestTokenMaker(t)

	token, _, err := tokenMaker.CreateAccessToken(ctx, contract.TokenPayloadInput{AccountUUID: accountUUID, SessionUUID: sessionUUID})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	req.Header.Set(infra.TokenKey.String(), token)
	return token
}

func GetTestContext(t *testing.T, req *http.Request, w http.ResponseWriter, authEndpoint bool) context.Context {
	t.Helper()

	c := echo.New().NewContext(req, w)
	if authEndpoint {
		c.Set(infra.AccountUUIDKey.String(), accountUUID)
		c.Set(infra.SessionKey.String(), sessionUUID)
	}
	return routeutils.GetContext(c)
}

type PrivateEndpointTest struct {
	Name          string
	Body          any
	SetupAuth     func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks)
	BuildMocks    func(ctx context.Context, m SvcMocks, body any)
	CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

var PrivateEndpointValidations = []PrivateEndpointTest{
	{
		Name: "Should return error when token is invalid",
		SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
			req.Header.Set(infra.TokenKey.String(), "invalid token")
		},
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			if recorder.Code != http.StatusUnauthorized {
				t.Errorf("Expected status code %d. Got %d", http.StatusUnauthorized, recorder.Code)
				t.Errorf("Response body: %s", recorder.Body.String())
			}
		},
	},
	{
		Name: "Should return error when token is missing",
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			if recorder.Code != http.StatusUnauthorized {
				t.Errorf("Expected status code %d. Got %d", http.StatusUnauthorized, recorder.Code)
				t.Errorf("Response body: %s", recorder.Body.String())
			}
		},
	},
	{
		Name: "Should return error when token is invalid",
		SetupAuth: func(ctx context.Context, t *testing.T, req *http.Request, m SvcMocks) {
			addAuthorizationWithNoCache(ctx, t, req)
		},
		BuildMocks: func(ctx context.Context, m SvcMocks, body any) {
			m.CacheMock.EXPECT().GetString(gomock.Any(), gomock.Any()).Return("invalid", nil).Times(1)
		},
		CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	},
}
