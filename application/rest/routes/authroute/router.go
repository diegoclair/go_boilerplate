package authroute

import (
	"github.com/labstack/echo/v4"
)

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

func (r *AccountRouter) RegisterRoutes(appGroup, privateGroup *echo.Group) {
	//TODO: create route for refresh token
	router := appGroup.Group(r.routeName)
	router.POST(loginRoute, r.ctrl.handleLogin)
}
