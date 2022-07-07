package authroute

import (
	"fmt"
	"sync"
	"time"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/domain/service"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/twinj/uuid"

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

	sessionUUID := uuid.NewV4().String()
	token, tokenPayload, err := s.authToken.CreateAccessToken(account.UUID, sessionUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	refreshToken, refreshTokenPayload, err := s.authToken.CreateRefreshToken(account.UUID, sessionUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	sessionReq := entity.Session{
		SessionUUID:           sessionUUID,
		AccountID:             account.ID,
		RefreshToken:          refreshToken,
		UserAgent:             c.Request().UserAgent(),
		ClientIP:              c.RealIP(),
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
	}

	err = s.authService.CreateSession(ctx, sessionReq)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.LoginResponse{
		AccessToken:           token,
		AccessTokenExpiresAt:  tokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	return routeutils.ResponseAPIOK(c, response)
}

func (s *Controller) handleRefreshToken(c echo.Context) error {

	ctx := routeutils.GetContext(c)

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
		return routeutils.HandleAPIError(c, err)
	}

	session, err := s.authService.GetSessionByUUID(ctx, refreshPayload.SessionUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	if session.IsBlocked {
		return routeutils.ResponseUnauthorizedError(c, fmt.Errorf("blocked session"))
	}
	if session.RefreshToken != input.RefreshToken {
		return routeutils.ResponseUnauthorizedError(c, fmt.Errorf("mismatched session token"))
	}
	if time.Now().After(session.RefreshTokenExpiredAt) {
		return routeutils.ResponseUnauthorizedError(c, fmt.Errorf("expired session"))
	}

	accessToken, accessPayload, err := s.authToken.CreateAccessToken(refreshPayload.AccountUUID, refreshPayload.SessionUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return routeutils.ResponseAPIOK(c, response)
}
