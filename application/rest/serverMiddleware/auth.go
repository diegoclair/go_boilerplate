package servermiddleware

import (
	"net/http"

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

			accessToken := ctx.Request().Header.Get(auth.ContextTokenKey.String())
			if len(accessToken) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, resterrors.NewUnauthorizedError("access token is required"))
			}

			payload, err := authToken.VerifyToken(accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			// Add information to the echo context
			ctx.Set(auth.AccountUUIDKey.String(), payload.AccountUUID)
			ctx.Set(auth.SessionKey.String(), payload.SessionUUID)

			return next(ctx)
		}
	}
}
