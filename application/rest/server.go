package rest

import (
	"fmt"

	"github.com/diegoclair/go-boilerplate/application/factory"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/accountroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/authroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/pingroute"
	"github.com/diegoclair/go-boilerplate/application/rest/routes/transferroute"
	servermiddleware "github.com/diegoclair/go-boilerplate/application/rest/serverMiddleware"
	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// IRouter interface for routers
type IRouter interface {
	RegisterRoutes(appGroup, privateGroup *echo.Group)
}

type Router struct {
	routers []IRouter
}

func StartRestServer(cfg *config.Config) {
	server := initServer(cfg)

	//TODO: create log package that we can pass sessionID and than we can trace user processes
	port := cfg.App.Port
	if port == "" {
		port = "5000"
	}

	log.Info(fmt.Sprintf("About to start the application on port: %s...", port))

	if err := server.Start(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}

func initServer(cfg *config.Config) *echo.Echo {

	factory := factory.GetDomainServices(cfg)

	srv := echo.New()
	srv.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	pingController := pingroute.NewController()
	accountController := accountroute.NewController(factory.AccountService, factory.Mapper)
	authController := authroute.NewController(factory.AuthService, factory.Mapper)
	transferController := transferroute.NewController(factory.TransferService, factory.Mapper)

	pingRoute := pingroute.NewRouter(pingController, "ping")
	accountRoute := accountroute.NewRouter(accountController, "accounts")
	authRoute := authroute.NewRouter(authController, "auth")
	transferRoute := transferroute.NewRouter(transferController, "transfers")

	appRouter := &Router{}
	appRouter.addRouters(accountRoute)
	appRouter.addRouters(authRoute)
	appRouter.addRouters(pingRoute)
	appRouter.addRouters(transferRoute)

	return appRouter.registerAppRouters(srv, cfg)
}

func (r *Router) addRouters(router IRouter) {
	r.routers = append(r.routers, router)
}

func (r *Router) registerAppRouters(srv *echo.Echo, cfg *config.Config) *echo.Echo {

	appGroup := srv.Group("/")
	privateGroup := appGroup.Group("",
		servermiddleware.JWTMiddlewareWithConfig(servermiddleware.JWTConfig{PrivateKey: cfg.App.Auth.JWTPrivateKey}),
		servermiddleware.JWTMiddlewarePrivateRoute())

	for _, appRouter := range r.routers {
		appRouter.RegisterRoutes(appGroup, privateGroup)
	}

	return srv
}
