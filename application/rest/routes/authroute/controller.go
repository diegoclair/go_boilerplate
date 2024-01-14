package authroute

import (
	"fmt"
	"sync"
	"time"

	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/twinj/uuid"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	authService contract.AuthService
	authToken   auth.AuthToken
	utils       routeutils.Utils
}

func NewController(authService contract.AuthService, authToken auth.AuthToken, utils routeutils.Utils) *Controller {
	once.Do(func() {
		instance = &Controller{
			authService: authService,
			authToken:   authToken,
			utils:       utils,
		}
	})
	return instance
}

func (s *Controller) handleLogin(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.Login{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}
	err = input.Validate()
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	account, err := s.authService.Login(ctx, input.CPF, input.Password)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	sessionUUID := uuid.NewV4().String()
	token, tokenPayload, err := s.authToken.CreateAccessToken(ctx, account.UUID, sessionUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	refreshToken, refreshTokenPayload, err := s.authToken.CreateRefreshToken(ctx, account.UUID, sessionUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
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
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := viewmodel.LoginResponse{
		AccessToken:           token,
		AccessTokenExpiresAt:  tokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	return s.utils.Resp().ResponseAPIOK(c, response)
}

func (s *Controller) handleRefreshToken(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.RefreshTokenRequest{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}
	err = input.Validate()
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	refreshPayload, err := s.authToken.VerifyToken(ctx, input.RefreshToken)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	session, err := s.authService.GetSessionByUUID(ctx, refreshPayload.SessionUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	if session.IsBlocked {
		return s.utils.Resp().ResponseUnauthorizedError(c, fmt.Errorf("blocked session"))
	}
	if session.RefreshToken != input.RefreshToken {
		return s.utils.Resp().ResponseUnauthorizedError(c, fmt.Errorf("mismatched session token"))
	}
	if time.Now().After(session.RefreshTokenExpiredAt) {
		return s.utils.Resp().ResponseUnauthorizedError(c, fmt.Errorf("expired session"))
	}

	accessToken, accessPayload, err := s.authToken.CreateAccessToken(ctx, refreshPayload.AccountUUID, refreshPayload.SessionUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := viewmodel.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return s.utils.Resp().ResponseAPIOK(c, response)
}
