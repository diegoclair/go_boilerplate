package authroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/goswag/models"
)

const RouteName = "auth"

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

func (r *AccountRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(r.routeName)

	router.POST(loginRoute, r.ctrl.handleLogin).
		Summary("Login").
		Read(viewmodel.Login{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.LoginResponse{},
			},
		})

	router.POST("/refresh-token", r.ctrl.handleRefreshToken).
		Summary("Refresh Token").
		Description("Generate a new token using the refresh token").
		Read(viewmodel.RefreshTokenRequest{}).
		Returns([]models.ReturnType{
			{
				StatusCode: http.StatusOK,
				Body:       viewmodel.RefreshTokenResponse{},
			},
		})
}
