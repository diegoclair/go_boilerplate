package pingroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/goswag/models"
)

const RouteName = "ping"

const (
	rootRoute = "/"
)

type PingRouter struct {
	ctrl      *Handler
	routeName string
}

func NewRouter(ctrl *Handler, routeName string) *PingRouter {
	return &PingRouter{
		ctrl:      ctrl,
		routeName: routeName,
	}
}

func (r *PingRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)

	router.GET(rootRoute, r.ctrl.handlePing).
		Summary("Ping the server").
		Description("Ping the server to check if it is alive").
		Returns([]models.ReturnType{
			{StatusCode: http.StatusOK, Body: pingResponse{}},
		})
}
