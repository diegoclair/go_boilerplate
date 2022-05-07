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

type Router struct {
	routers []IRouter
}

func StartRestServer(cfg *config.Config, services *factory.Services, log logger.Logger, authToken auth.AuthToken) {
	server := initServer(cfg, services, authToken)
	port := cfg.App.Port
	if port == "" {
		port = "5000"
	}

	log.Info(fmt.Sprintf("About to start the application on port: %s...", port))

	if err := server.Start(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}

func initServer(cfg *config.Config, services *factory.Services, authToken auth.AuthToken) *echo.Echo {

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

	appRouter := &Router{}
	appRouter.addRouters(accountRoute)
	appRouter.addRouters(authRoute)
	appRouter.addRouters(pingRoute)
	appRouter.addRouters(transferRoute)

	return appRouter.registerAppRouters(srv, cfg, authToken)
}

func (r *Router) addRouters(router IRouter) {
	r.routers = append(r.routers, router)
}

func (r *Router) registerAppRouters(srv *echo.Echo, cfg *config.Config, authToken auth.AuthToken) *echo.Echo {

	appGroup := srv.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.AuthMiddlewarePrivateRoute(authToken))

	for _, appRouter := range r.routers {
		appRouter.RegisterRoutes(appGroup, privateGroup)
	}

	return srv
}
