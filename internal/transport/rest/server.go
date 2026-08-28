package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	infraContract "github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/internal/application/service"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/clientip"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/pingroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/swaggerroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/internal/transport/rest/serverMiddleware"
	"github.com/diegoclair/goswag/v2"
	"github.com/diegoclair/logger"
	echoprometheus "github.com/labstack/echo-prometheus"
	echo "github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	defaultPort     = "5000"
	gracefulTimeout = 10 * time.Second
)

type Server struct {
	routes []routeutils.IRoute
	Router goswag.Echo
	cache  contract.CacheManager
	log    logger.Logger
	port   string
}

func NewServer(cfg *config.Config, infra domain.Infrastructure, services *service.Apps, appName, port string) *Server {
	server := NewRestServer(services, cfg.GetAuthToken(), infra.CacheManager(), appName, cfg.App.ClientIPHeaders)
	if port == "" {
		port = defaultPort
	}

	server.port = port
	server.log = infra.Logger()

	return server
}

func NewRestServer(services *service.Apps, authToken infraContract.AuthToken, cache contract.CacheManager, appName string, clientIPHeaders []string) *Server {
	router := goswag.NewEcho(routeutils.DefaultSwaggerErrors()...)
	router.Echo().IPExtractor = clientip.Extractor(clientIPHeaders, nil)
	router.Echo().Use(middleware.CORS("*"))
	router.Echo().HTTPErrorHandler = func(c *echo.Context, err error) {
		_ = routeutils.HandleError(c, err)
	}

	pingHandler := pingroute.NewHandler()
	accountHandler := accountroute.NewHandler(services.AccountService)
	authHandler := authroute.NewHandler(services.AuthService, authToken)
	transferHandler := transferroute.NewHandler(services.TransferService)

	pingRoute := pingroute.NewRouter(pingHandler)
	accountRoute := accountroute.NewRouter(accountHandler)
	authRoute := authroute.NewRouter(authHandler)
	transferRoute := transferroute.NewRouter(transferHandler)

	swaggerRoute := swaggerroute.NewRouter(router.Echo())

	server := &Server{Router: router, cache: cache}
	server.addRouters(accountRoute)
	server.addRouters(authRoute)
	server.addRouters(pingRoute)
	server.addRouters(transferRoute)
	server.addRouters(swaggerRoute)
	server.registerAppRouters(authToken)

	server.setupPrometheus(appName)

	return server
}

func (r *Server) addRouters(router routeutils.IRoute) {
	r.routes = append(r.routes, router)
}

func (r *Server) registerAppRouters(authToken infraContract.AuthToken) {
	g := &routeutils.EchoGroups{}
	g.AppGroup = r.Router.Group("/")
	g.PrivateGroup = g.AppGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken, r.cache),
	)

	for _, appRouter := range r.routes {
		appRouter.RegisterRoutes(g)
	}
}

func (r *Server) setupPrometheus(appName string) {
	p := echoprometheus.NewMiddleware(appName)
	r.Router.Echo().Use(p)
}

// Serve blocks until ctx is cancelled, then drains the requests in flight —
// bounded by gracefulTimeout. Nothing the requests depend on may close before
// it returns.
func (r *Server) Serve(ctx context.Context) error {
	r.log.Info(ctx, fmt.Sprintf("About to start the application on port: %s...", r.port))

	return echo.StartConfig{
		Address:         fmt.Sprintf(":%s", r.port),
		GracefulTimeout: gracefulTimeout,
		OnShutdownError: func(err error) {
			r.log.Error(ctx, "Failed to shutdown rest server", logger.Err(err))
		},
	}.Start(ctx, r.Router.Echo())
}
