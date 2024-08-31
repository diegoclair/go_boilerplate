package servermiddleware

import (
	"net/http"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_utils/resterrors"
	echo "github.com/labstack/echo/v4"
)

func AuthMiddlewarePrivateRoute(authToken auth.AuthToken, cache contract.CacheManager) echo.MiddlewareFunc {
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

			valid, _ := cache.GetString(ctx.Request().Context(), accessToken)
			if valid != "" {
				return echo.NewHTTPError(http.StatusUnauthorized, resterrors.NewUnauthorizedError("token is invalid"))
			}

			// Add information to the echo context
			ctx.Set(infra.AccountUUIDKey.String(), payload.AccountUUID)
			ctx.Set(infra.SessionKey.String(), payload.SessionUUID)

			return next(ctx)
		}
	}
}
