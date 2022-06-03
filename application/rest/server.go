package rest

import (
	"fmt"

	"github.com/diegoclair/go-boilerplate/application/factory"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/accountroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/authroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/pingroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/transferroute"
	servermiddleware "github.com/diegoclair/go-boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// IRouter interface for routers
type IRouter interface {
	RegisterRoutes(appGroup, privateGroup *echo.Group)
}

type Server struct {
	routers []IRouter
	server  *echo.Echo
}

func StartRestServer(cfg *config.Config, services *factory.Services, log logger.Logger, authToken auth.AuthToken) {
	server := newRestServer(services, authToken)
	port := cfg.App.Port
	if port == "" {
		port = "5000"
	}

	log.Info(fmt.Sprintf("About to start the application on port: %s...", port))

	if err := server.Start(port); err != nil {
		panic(err)
	}
}

func newRestServer(services *factory.Services, authToken auth.AuthToken) *Server {

	srv := echo.New()
	srv.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	pingController := pingroute.NewController()
	accountController := accountroute.NewController(services.AccountService, services.Mapper)
	authController := authroute.NewController(services.AuthService, services.Mapper, authToken)
	transferController := transferroute.NewController(services.TransferService, services.Mapper)

	pingRoute := pingroute.NewRouter(pingController, "ping")
	accountRoute := accountroute.NewRouter(accountController, "accounts")
	authRoute := authroute.NewRouter(authController, "auth")
	transferRoute := transferroute.NewRouter(transferController, "transfers")

	server := &Server{server: srv}
	server.addRouters(accountRoute)
	server.addRouters(authRoute)
	server.addRouters(pingRoute)
	server.addRouters(transferRoute)
	server.registerAppRouters(authToken)

	return server
}

func (r *Server) addRouters(router IRouter) {
	r.routers = append(r.routers, router)
}

func (r *Server) registerAppRouters(authToken auth.AuthToken) {

	appGroup := r.server.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken))

	for _, appRouter := range r.routers {
		appRouter.RegisterRoutes(appGroup, privateGroup)
	}

	return
}

func (r *Server) Start(port string) error {
	return r.server.Start(fmt.Sprintf(":%s", port))
}
