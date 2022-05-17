package authroute

import (
	"sync"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go-boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go-boilerplate/domain/service"
	"github.com/diegoclair/go-boilerplate/infra/auth"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	authService service.AuthService
	mapper      mapper.Mapper
	authToken   auth.AuthToken
}

func NewController(authService service.AuthService, mapper mapper.Mapper, authToken auth.AuthToken) *Controller {
	once.Do(func() {
		instance = &Controller{
			authService: authService,
			mapper:      mapper,
			authToken:   authToken,
		}
	})
	return instance
}

func (s *Controller) handleLogin(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	input := viewmodel.Login{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}
	err = input.Validate()
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	account, err := s.authService.Login(ctx, input.CPF, input.Secret)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	token, tokenPayload, err := s.authToken.CreateToken(account.UUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	refreshToken, refreshTokenPayload, err := s.authToken.CreateRefreshToken(account.UUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.LoginResponse{
		AccessToken:           token,
		AccessTokenExpiresAt:  tokenPayload.ExpiredAt.Unix(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt.Unix(),
	}

	return routeutils.ResponseAPIOK(c, response)
}

func (s *Controller) handleRefreshToken(c echo.Context) error {

	input := viewmodel.RefreshTokenRequest{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}
	err = input.Validate()
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	refreshPayload, err := s.authToken.VerifyToken(input.RefreshToken)
	if err != nil {
		return routeutils.ResponseUnauthorizedError(c, err)
	}

	//TODO: handle session here
	// session, err := s.store.GetSession(ctx, refreshPayload.ID)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }
	// if session.IsBlocked {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("blocked session")))
	// 	return
	// }
	// if session.Username != refreshPayload.Username {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect session user")))
	// 	return
	// }

	// if session.RefreshToken != input.RefreshToken {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("mismatched session token")))
	// 	return
	// }
	// if time.Now().After(session.ExpiresAt) {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired session")))
	// 	return
	// }

	accessToken, accessPayload, err := s.authToken.CreateToken(refreshPayload.AccountUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return routeutils.ResponseAPIOK(c, response)
}
