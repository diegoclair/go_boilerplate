package shared

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/mocks"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/transport/rest/serverMiddleware"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/goswag"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"go.uber.org/mock/gomock"
)

type SvcMocks struct {
	AccountMock   *mocks.MockAccountService
	AuthMock      *mocks.MockAuthService
	AuthTokenMock *mocks.MockAuthToken
	TransferMock  *mocks.MockTransferService
}

func GetServerTest(t *testing.T) (m SvcMocks, server goswag.Echo, ctrl *gomock.Controller) {
	t.Helper()

	ctrl = gomock.NewController(t)
	m = SvcMocks{
		AccountMock:   mocks.NewMockAccountService(ctrl),
		AuthMock:      mocks.NewMockAuthService(ctrl),
		AuthTokenMock: mocks.NewMockAuthToken(ctrl),
		TransferMock:  mocks.NewMockTransferService(ctrl),
	}

	server = goswag.NewEcho()
	appGroup := server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(getTestTokenMaker(t)),
	)

	g := &routeutils.EchoGroups{
		AppGroup:     appGroup,
		PrivateGroup: privateGroup,
	}

	accountHandler := accountroute.NewHandler(m.AccountMock)
	accountRoute := accountroute.NewRouter(accountHandler, accountroute.RouteName)
	authHandler := authroute.NewHandler(m.AuthMock, m.AuthTokenMock)
	authRoute := authroute.NewRouter(authHandler, authroute.RouteName)
	transferHandler := transferroute.NewHandler(m.TransferMock)
	transferRoute := transferroute.NewRouter(transferHandler, transferroute.RouteName)

	accountRoute.RegisterRoutes(g)
	authRoute.RegisterRoutes(g)
	transferRoute.RegisterRoutes(g)
	return
}

var (
	tokenMaker auth.AuthToken
	onceToken  sync.Once
)

func getTestTokenMaker(t *testing.T) auth.AuthToken {
	t.Helper()

	onceToken.Do(func() {
		cfg, err := config.GetConfigEnvironment(config.ProfileTest)
		require.NoError(t, err)

		cfg.App.Auth.AccessTokenDuration = 2 * time.Second
		cfg.App.Auth.RefreshTokenDuration = 2 * time.Second

		tokenMaker, err = auth.NewAuthToken(cfg.App.Auth, logger.NewNoop())
		require.NoError(t, err)
	})
	return tokenMaker
}

var (
	accountUUID = uuid.NewV4().String()
	sessionUUID = uuid.NewV4().String()
)

func AddAuthorization(ctx context.Context, t *testing.T, req *http.Request) {
	t.Helper()

	tokenMaker := getTestTokenMaker(t)

	token, _, err := tokenMaker.CreateAccessToken(ctx, auth.TokenPayloadInput{AccountUUID: accountUUID, SessionUUID: sessionUUID})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	req.Header.Set(infra.TokenKey.String(), token)
}

func GetTestContext(t *testing.T, req *http.Request, w http.ResponseWriter) context.Context {
	t.Helper()

	c := echo.New().NewContext(req, w)
	c.Set(infra.AccountUUIDKey.String(), accountUUID)
	c.Set(infra.SessionKey.String(), sessionUUID)
	return routeutils.GetContext(c)
}
