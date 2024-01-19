package transferroute

import "github.com/diegoclair/go_boilerplate/application/rest/routeutils"

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

func (r *TransferRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.PrivateGroup.Group(r.routeName)
	router.POST(rootRoute, r.ctrl.handleAddTransfer)
	router.GET(rootRoute, r.ctrl.handleGetTransfers)
}
