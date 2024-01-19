package pingroute

import (
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
)

const RouteName = "ping"

const (
	rootRoute = "/"
)

type PingRouter struct {
	ctrl      *Controller
	routeName string
}

func NewRouter(ctrl *Controller, routeName string) *PingRouter {
	return &PingRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *PingRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)
	router.GET(rootRoute, r.ctrl.handlePing)
}
