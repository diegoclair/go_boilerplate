package routeutils

import (
	"context"
	"strconv"
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

// GetPagingParams gets the standard paging params from the URL, returning how much data to take and skip
func GetPagingParams(c echo.Context, pageParameter, quantityParameter string) (take int64, skip int64) {

	if pageParameter == "" {
		pageParameter = "page"
	}

	if quantityParameter == "" {
		quantityParameter = "quantity"
	}
	pg := c.QueryParam(pageParameter)
	ipp := c.QueryParam(quantityParameter)

	page, _ := strconv.ParseInt(pg, 10, 64)
	quantity, _ := strconv.ParseInt(ipp, 10, 64)

	if page < 1 {
		page = 1
	}

	if quantity < 1 || quantity > 1000 {
		quantity = 10
	}

	take = quantity // itens per page
	skip = (page - 1) * quantity

	return
}
