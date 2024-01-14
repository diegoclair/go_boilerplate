package rest

import (
	"context"
	"fmt"

	"github.com/diegoclair/go_boilerplate/application/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/application/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/application/rest/routes/pingroute"
	"github.com/diegoclair/go_boilerplate/application/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go_boilerplate/domain/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_utils-lib/v2/validator"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// IRouter interface for routers
type IRouter interface {
	RegisterRoutes(appGroup, privateGroup *echo.Group)
}

type Server struct {
	routers []IRouter
	Srv     *echo.Echo
	cfg     *config.Config
}

func StartRestServer(ctx context.Context, cfg *config.Config, services *service.Services, log logger.Logger, authToken auth.AuthToken, v validator.Validator) {
	server := newRestServer(services, authToken, cfg, log, v)
	port := cfg.App.Port
	if port == "" {
		port = "5000"
	}

	log.Infof(ctx, "About to start the application on port: %s...", port)

	if err := server.Start(port); err != nil {
		panic(err)
	}
}

func newRestServer(services *service.Services, authToken auth.AuthToken, cfg *config.Config, log logger.Logger, v validator.Validator) *Server {

	srv := echo.New()
	srv.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	routeUtils := routeutils.New(log)

	pingController := pingroute.NewController()
	accountController := accountroute.NewController(services.AccountService, routeUtils, v)
	authController := authroute.NewController(services.AuthService, authToken, routeUtils, v)
	transferController := transferroute.NewController(services.TransferService, routeUtils, v)

	pingRoute := pingroute.NewRouter(pingController, pingroute.RouteName)
	accountRoute := accountroute.NewRouter(accountController, accountroute.RouteName)
	authRoute := authroute.NewRouter(authController, authroute.RouteName)
	transferRoute := transferroute.NewRouter(transferController, transferroute.RouteName)

	server := &Server{Srv: srv, cfg: cfg}
	server.addRouters(accountRoute)
	server.addRouters(authRoute)
	server.addRouters(pingRoute)
	server.addRouters(transferRoute)
	server.registerAppRouters(authToken)

	server.setupPrometheus()

	return server
}

func (r *Server) addRouters(router IRouter) {
	r.routers = append(r.routers, router)
}

func (r *Server) registerAppRouters(authToken auth.AuthToken) {

	appGroup := r.Srv.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken),
	)

	for _, appRouter := range r.routers {
		appRouter.RegisterRoutes(appGroup, privateGroup)
	}

}

func (r *Server) setupPrometheus() {

	p := prometheus.NewPrometheus(r.cfg.App.Name, nil)
	p.Use(r.Srv)
}

func (r *Server) Start(port string) error {
	return r.Srv.Start(fmt.Sprintf(":%s", port))
}
