package swaggerroute

import (
	_ "github.com/diegoclair/go_boilerplate/docs"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	echo "github.com/labstack/echo/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type swaggerRouter struct {
	e *echo.Echo
}

func NewRouter(e *echo.Echo) *swaggerRouter {
	return &swaggerRouter{e: e}
}

// TODO: this route should not be public, only internal (maybe staging or dev environments)
func (r *swaggerRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := r.e.Group("swagger")
	router.GET("/*", echo.WrapHandler(httpSwagger.Handler()))
}
