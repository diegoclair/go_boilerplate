package transferroute

import (
	"github.com/labstack/echo/v4"
)

const RouteName = "transfers"

const (
	rootRoute = ""
)

type TransferRouter struct {
	ctrl      *Controller
	routeName string
}

func NewRouter(ctrl *Controller, routeName string) *TransferRouter {
	return &TransferRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *TransferRouter) RegisterRoutes(appGroup, privateGroup *echo.Group) {
	router := privateGroup.Group(r.routeName)
	router.POST(rootRoute, r.ctrl.handleAddTransfer)
	router.GET(rootRoute, r.ctrl.handleGetTransfers)
}
