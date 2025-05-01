package transferroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/viewmodel"
	"github.com/diegoclair/goswag"
	"github.com/diegoclair/goswag/models"
)

const GroupRouteName = "transfers"

const (
	RootRoute = ""
)

type TransferRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *TransferRouter {
	return &TransferRouter{
		ctrl: ctrl,
	}
}

func (r *TransferRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.PrivateGroup.Group(GroupRouteName)

	router.POST(RootRoute, r.ctrl.handleAddTransfer).
		Summary("Add a new transfer").
		Read(viewmodel.TransferReq{}).
		Returns([]models.ReturnType{{StatusCode: http.StatusCreated}}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)

	router.GET(RootRoute, r.ctrl.handleGetTransfers).
		Summary("Get all transfers").
		Description("Get all transfers with paginated response").
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.PaginatedResponse[[]viewmodel.TransferResp]{},
			},
		}).
		HeaderParam(infra.TokenKey.String(), infra.TokenKeyDescription, goswag.StringType, true)
}
