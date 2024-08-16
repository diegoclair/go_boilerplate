package transferroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/goswag/models"
)

const RouteName = "transfers"

const (
	RootRoute = ""
)

type TransferRouter struct {
	ctrl      *Handler
	routeName string
}

func NewRouter(ctrl *Handler, routeName string) *TransferRouter {
	return &TransferRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *TransferRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.PrivateGroup.Group(r.routeName)

	router.POST(RootRoute, r.ctrl.handleAddTransfer).
		Summary("Add a new transfer").
		Read(viewmodel.TransferReq{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}})

	router.GET(RootRoute, r.ctrl.handleGetTransfers).
		Summary("Get all transfers").
		Description("Get all transfers with paginated response").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResponse[[]viewmodel.TransferResp]{},
			},
		})
}
