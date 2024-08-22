package authroute

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/goswag/models"
)

const GroupRouteName = "auth"

const (
	LoginRoute = "/login"
)

type AuthRouter struct {
	ctrl *Handler
}

func NewRouter(ctrl *Handler) *AuthRouter {
	return &AuthRouter{
		ctrl: ctrl,
	}
}

func (r *AuthRouter) RegisterRoutes(g *routeutils.EchoGroups) {
	router := g.AppGroup.Group(GroupRouteName)

	router.POST(LoginRoute, r.ctrl.handleLogin).
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
