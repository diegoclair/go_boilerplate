package authroute

import (
	"context"
	"sync"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	infraContract "github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/viewmodel"
	"github.com/twinj/uuid"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	authService contract.AuthApp
	authToken   infraContract.AuthToken
}

func NewHandler(authService contract.AuthApp, authToken infraContract.AuthToken) *Handler {
	Once.Do(func() {
		instance = &Handler{
			authService: authService,
			authToken:   authToken,
		}
	})

	return instance
}

func (s *Handler) handleLogin(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.Login{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	account, err := s.authService.Login(ctx, input.ToDto())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	sessionUUID := uuid.NewV4().String()
	req := infraContract.TokenPayloadInput{
		AccountUUID: account.UUID,
		SessionUUID: sessionUUID,
	}
	token, tokenPayload, err := s.authToken.CreateAccessToken(ctx, req)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	refreshToken, refreshTokenPayload, err := s.authToken.CreateRefreshToken(ctx, req)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	sessionReq := dto.Session{
		SessionUUID:           sessionUUID,
		AccountID:             account.ID,
		RefreshToken:          refreshToken,
		UserAgent:             c.Request().UserAgent(),
		ClientIP:              c.RealIP(),
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
	}

	err = s.authService.CreateSession(ctx, sessionReq)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.LoginResponse{
		AccessToken:           token,
		AccessTokenExpiresAt:  tokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleRefreshToken(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.RefreshTokenRequest{}

	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	refreshPayload, err := s.authToken.VerifyToken(ctx, input.RefreshToken)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	ctx = context.WithValue(ctx, infra.AccountUUIDKey, refreshPayload.AccountUUID)
	ctx = context.WithValue(ctx, infra.SessionKey, refreshPayload.SessionUUID)

	session, err := s.authService.GetSessionByUUID(ctx, refreshPayload.SessionUUID)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	if session.IsBlocked {
		return routeutils.ResponseUnauthorizedError(c, "session blocked")
	}

	if session.RefreshToken != input.RefreshToken {
		return routeutils.ResponseUnauthorizedError(c, "mismatched session token")
	}

	if time.Now().After(session.RefreshTokenExpiredAt) {
		return routeutils.ResponseUnauthorizedError(c, "expired session")
	}

	req := infraContract.TokenPayloadInput{
		AccountUUID: refreshPayload.AccountUUID,
		SessionUUID: refreshPayload.SessionUUID,
	}
	accessToken, accessPayload, err := s.authToken.CreateAccessToken(ctx, req)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := viewmodel.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return routeutils.ResponseAPIOk(c, response)
}

func (s *Handler) handleLogout(c echo.Context) error {
	accessToken := c.Request().Header.Get(infra.TokenKey.String())
	ctx := routeutils.GetContext(c)

	err := s.authService.Logout(ctx, accessToken)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseAPIOk(c, struct{}{})
}
