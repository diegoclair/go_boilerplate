package servermiddleware

import (
	"fmt"

	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/labstack/echo/v4"
)

// JWTConfig defines the config for JWT middleware.
type JWTConfig struct {
	PrivateKey string
}

func AuthMiddlewarePrivateRoute(authToken auth.AuthToken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			token := ctx.Get(auth.ContextTokenKey.String())

			//TODO: finish token middleware
			fmt.Println(token)

			// if !token.Valid {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// }

			// claims, ok := token.Claims.(*entity.TokenData)
			// if !ok {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// }

			// if !claims.LoggedIn {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// }

			// Add account information to context
			//ctx.Set(auth.AccountUUIDKey.String(), claims.AccountUUID)

			return next(ctx)
		}
	}
}
