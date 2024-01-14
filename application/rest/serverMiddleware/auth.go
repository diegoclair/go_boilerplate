package servermiddleware

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
)

// JWTConfig defines the config for JWT middleware.
type JWTConfig struct {
	PrivateKey string
}

func AuthMiddlewarePrivateRoute(authToken auth.AuthToken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			accessToken := ctx.Request().Header.Get(infra.TokenKey.String())
			if len(accessToken) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, resterrors.NewUnauthorizedError("access token is required"))
			}

			payload, err := authToken.VerifyToken(ctx.Request().Context(), accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			// Add information to the echo context
			ctx.Set(infra.AccountUUIDKey.String(), payload.AccountUUID)
			ctx.Set(infra.SessionKey.String(), payload.SessionUUID)

			return next(ctx)
		}
	}
}
