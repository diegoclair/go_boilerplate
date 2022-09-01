package accountroute

import (
	"github.com/labstack/echo/v4"
)

const RouteName = "accounts"

const (
	rootRoute          = ""
	accountByID        = "/:account_uuid"
	accountBalanceByID = "/:account_uuid/balance"
)

type UserRouter struct {
	ctrl      *Controller
	routeName string
}

func NewRouter(ctrl *Controller, routeName string) *UserRouter {
	return &UserRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *UserRouter) RegisterRoutes(appGroup, privateGroup *echo.Group) {
	router := appGroup.Group(r.routeName)
	router.POST(rootRoute, r.ctrl.handleAddAccount)
	router.POST(accountBalanceByID, r.ctrl.handleAddBalance)
	router.GET(rootRoute, r.ctrl.handleGetAccounts)
	router.GET(accountByID, r.ctrl.handleGetAccountByID)
	router.GET(accountBalanceByID, r.ctrl.handleGetAccountBalanceByID)

}
