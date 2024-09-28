package swaggerroute

import (
	_ "github.com/diegoclair/go_boilerplate/docs"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type swaggerRouter struct {
}

func NewRouter() *swaggerRouter {
	return &swaggerRouter{}
}

// TODO: this route should not be public, only internal (maybe staging or dev environments)
func (r *swaggerRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group("swagger")
	router.GET("/*", echoSwagger.WrapHandler)
}
