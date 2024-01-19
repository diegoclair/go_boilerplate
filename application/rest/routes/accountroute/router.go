package accountroute

import (
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
)

const RouteName = "accounts"

const (
	rootRoute          = ""
	accountByID        = "/:account_uuid/"
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

func (r *UserRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)
	router.POST(rootRoute, r.ctrl.handleAddAccount)
	router.POST(accountBalanceByID, r.ctrl.handleAddBalance)
	router.GET(rootRoute, r.ctrl.handleGetAccounts)
	router.GET(accountByID, r.ctrl.handleGetAccountByID)
	router.GET(accountBalanceByID, r.ctrl.handleGetAccountBalanceByID)
}
