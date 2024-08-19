package authroute

import (
	"sync"
	"time"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/twinj/uuid"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	authService contract.AuthService
	authToken   auth.AuthToken
}

func NewHandler(authService contract.AuthService, authToken auth.AuthToken) *Handler {
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
		return routeutils.ResponseInvalidRequestBody(c)
	}

	account, err := s.authService.Login(ctx, input.ToDto())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	sessionUUID := uuid.NewV4().String()
	req := auth.TokenPayloadInput{
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
		return routeutils.ResponseInvalidRequestBody(c)
	}

	refreshPayload, err := s.authToken.VerifyToken(ctx, input.RefreshToken)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

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

	req := auth.TokenPayloadInput{
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
