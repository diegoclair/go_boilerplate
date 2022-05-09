package routeutils

import (
	"context"
	"strings"

	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
)

// GetContext returns a fulled appcontext
func GetContext(ctx echo.Context) (appContext context.Context) {
	appContext = context.Background()
	appContext = context.WithValue(appContext, auth.AccountUUIDKey, ctx.Get(auth.AccountUUIDKey.String()))
	appContext = context.WithValue(appContext, auth.SessionKey, ctx.Get(auth.SessionKey.String()))
	return appContext
}

// GetAndValidateParam gets the param value and validates it, returning a validation error in case it's invalid
func GetAndValidateParam(c echo.Context, paramName string, errorMessage string) (paramValue string, err error) {
	paramValue = c.Param(paramName)

	if strings.TrimSpace(paramValue) == "" {
		return paramValue, resterrors.NewUnprocessableEntity(errorMessage)
	}

	return paramValue, nil
}
