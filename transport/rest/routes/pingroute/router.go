package pingroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/goswag/models"
)

const GroupRouteName = "ping"

const (
	rootRoute = "/"
)

type PingRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *PingRouter {
	return &PingRouter{
		ctrl: ctrl,
	}
}

func (r *PingRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(GroupRouteName)

	router.GET(rootRoute, r.ctrl.handlePing).
		Summary("Ping the server").
		Description("Ping the server to check if it is alive").
		Returns([]models.ReturnType{
			{StatusCode: http.StatusOK, Body: pingResponse{}},
		})
}
