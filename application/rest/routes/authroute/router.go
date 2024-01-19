package authroute

import "github.com/diegoclair/go_boilerplate/application/rest/routeutils"

const RouteName = "auth"

const (
	loginRoute = "/login"
)

type AccountRouter struct {
	ctrl      *Controller
	routeName string
}

func NewRouter(ctrl *Controller, routeName string) *AccountRouter {
	return &AccountRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *AccountRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)
	router.POST(loginRoute, r.ctrl.handleLogin)
	router.POST("/refresh-token", r.ctrl.handleRefreshToken)

}
