package routeutils

import (
	"context"
	"strconv"
	"strings"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_utils/resterrors"
	echo "github.com/labstack/echo/v4"
)

func GetContext(c echo.Context) (ctx context.Context) {
	ctx = c.Request().Context()
	ctx = context.WithValue(ctx, infra.AccountUUIDKey, c.Get(infra.AccountUUIDKey.String()))
	ctx = context.WithValue(ctx, infra.SessionKey, c.Get(infra.SessionKey.String()))
	return ctx
}

// GetRequiredInt64PathParam gets the param value and validates it, returning a validation error in case it's empty
func GetRequiredInt64PathParam(c echo.Context, paramName string, errorMessage string) (paramValue int64, err error) {
	paramValue, err = strconv.ParseInt(c.Param(paramName), 10, 64)
	if err != nil {
		return paramValue, resterrors.NewUnprocessableEntity(errorMessage)
	}

	if paramValue == 0 {
		return paramValue, resterrors.NewUnprocessableEntity(errorMessage)
	}

	return paramValue, nil
}

// GetRequiredStringPathParam gets the param value and validates it, returning a validation error in case it's empty
func GetRequiredStringPathParam(c echo.Context, paramName string, errorMessage string) (paramValue string, err error) {
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

	return GetTakeSkipFromPageQuantity(page, quantity)
}

func GetTakeSkipFromPageQuantity(page, quantity int64) (take, skip int64) {
	if page < 1 {
		page = 1
	}

	if quantity < 1 || quantity > 1000 {
		quantity = 10
	}

	take = quantity // items per page
	skip = (page - 1) * quantity
	return
}
