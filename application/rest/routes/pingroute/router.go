package pingroute

import (
	"github.com/labstack/echo/v4"
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

func (r *PingRouter) RegisterRoutes(appGroup, privateGroup *echo.Group) {
	router := appGroup.Group(r.routeName)
	router.GET(rootRoute, r.ctrl.handlePing)
}
