package servermiddleware

import (
	"github.com/diegoclair/go_boilerplate/infra"
	infraContract "github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_utils/resterrors"
	echo "github.com/labstack/echo/v4"
)

func AuthMiddlewarePrivateRoute(authToken infraContract.AuthToken, cache contract.CacheManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			accessToken := ctx.Request().Header.Get(infra.TokenKey.String())
			if len(accessToken) == 0 {
				return resterrors.NewUnauthorizedError("access token is required")
			}

			payload, err := authToken.VerifyToken(ctx.Request().Context(), accessToken)
			if err != nil {
				return err
			}

			valid, _ := cache.GetString(ctx.Request().Context(), accessToken)
			if valid != "" {
				return resterrors.NewUnauthorizedError("token is invalid")
			}

			// Add information to the echo context
			ctx.Set(infra.AccountUUIDKey.String(), payload.AccountUUID)
			ctx.Set(infra.SessionKey.String(), payload.SessionUUID)

			return next(ctx)
		}
	}
}
