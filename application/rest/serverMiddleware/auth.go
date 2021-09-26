package servermiddleware

import (
	"net/http"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTConfig defines the config for JWT middleware.
type JWTConfig struct {
	PrivateKey string
}

// JWTMiddlewareWithConfig returns a JWT middleware with config.
func JWTMiddlewareWithConfig(jwtConfig JWTConfig) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		TokenLookup:   "header:" + auth.TokenHeaderName.String(),
		ContextKey:    auth.ContextTokenKey.String(),
		SigningKey:    []byte(jwtConfig.PrivateKey),
		SigningMethod: auth.TokenSigningMethod.Name,
		Claims:        &entity.TokenData{},
	})
}

func JWTMiddlewarePrivateRoute() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			token, ok := ctx.Get(auth.ContextTokenKey.String()).(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			claims, ok := token.Claims.(*entity.TokenData)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			if !claims.LoggedIn {
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			// Add account information to context
			ctx.Set(auth.AccountUUIDKey.String(), claims.AccountUUID)

			return next(ctx)
		}
	}
}
