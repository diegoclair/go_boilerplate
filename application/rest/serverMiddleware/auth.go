package servermiddleware

import (
	"net/http"

	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			payload, err := authToken.VerifyToken(accessToken)
			if err != nil {
				log.Error("Erro")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			//TODO: add session here too
			// Add information to context
			ctx.Set(auth.AccountUUIDKey.String(), payload.AccountUUID)

			return next(ctx)
		}
	}
}
