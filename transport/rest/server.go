package rest

import (
	"context"
	"fmt"

	"github.com/diegoclair/go_boilerplate/application/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/accountroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/authroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/pingroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routes/transferroute"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	servermiddleware "github.com/diegoclair/go_boilerplate/transport/rest/serverMiddleware"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/goswag"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	routes []routeutils.IRoute
	Router goswag.Echo
	cfg    *config.Config
}

func StartRestServer(ctx context.Context, cfg *config.Config, services *service.Services, log logger.Logger, authToken auth.AuthToken) *Server {
	server := NewRestServer(services, authToken, cfg)
	port := cfg.App.Port
	if port == "" {
		port = "5000"
	}

	log.Infof(ctx, "About to start the application on port: %s...", port)

	if err := server.Start(port); err != nil {
		panic(err)
	}

	go func() {
		if err := server.Start(port); err != nil {
			panic(err)
		}
	}()

	return server
}

func NewRestServer(services *service.Services, authToken auth.AuthToken, cfg *config.Config) *Server {
	router := goswag.NewEcho()
	router.Echo().Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	pingHandler := pingroute.NewHandler()
	accountHandler := accountroute.NewHandler(services.AccountService)
	authHandler := authroute.NewHandler(services.AuthService, authToken)
	transferHandler := transferroute.NewHandler(services.TransferService)

	pingRoute := pingroute.NewRouter(pingHandler, pingroute.RouteName)
	accountRoute := accountroute.NewRouter(accountHandler, accountroute.RouteName)
	authRoute := authroute.NewRouter(authHandler, authroute.RouteName)
	transferRoute := transferroute.NewRouter(transferHandler, transferroute.RouteName)

	server := &Server{Router: router, cfg: cfg}
	server.addRouters(accountRoute)
	server.addRouters(authRoute)
	server.addRouters(pingRoute)
	server.addRouters(transferRoute)
	server.registerAppRouters(authToken)

	server.setupPrometheus()

	return server
}

func (r *Server) addRouters(router routeutils.IRoute) {
	r.routes = append(r.routes, router)
}

func (r *Server) registerAppRouters(authToken auth.AuthToken) {
	g := &routeutils.EchoGroups{}
	g.AppGroup = r.Router.Group("/")
	g.PrivateGroup = g.AppGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken),
	)

	for _, appRouter := range r.routes {
		appRouter.RegisterRoutes(g)
	}
}

func (r *Server) setupPrometheus() {
	p := prometheus.NewPrometheus(r.cfg.App.Name, nil)
	p.Use(r.Router.Echo())
}

func (r *Server) Start(port string) error {
	return r.Router.Echo().Start(fmt.Sprintf(":%s", port))
}
