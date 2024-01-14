package routeutils

import (
	"context"
	"strconv"
	"strings"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
)

type RequestUtils interface {
	// GetContext returns a filled ctx with the account uuid and session code if route has access token
	GetContext(c echo.Context) (ctx context.Context)
	GetAndValidateParam(c echo.Context, paramName string, errorMessage string) (paramValue string, err error)
	GetPagingParams(c echo.Context, pageParameter, quantityParameter string) (take int64, skip int64)
	GetTakeSkipFromPageQuantity(page, quantity int64) (take, skip int64)
}

type reqUtils struct{}

func newRequestUtils() RequestUtils {
	return &reqUtils{}
}

// GetContext returns a filled ctx
func (r *reqUtils) GetContext(c echo.Context) (ctx context.Context) {
	ctx = c.Request().Context()
	ctx = context.WithValue(ctx, infra.AccountUUIDKey, c.Get(infra.AccountUUIDKey.String()))
	ctx = context.WithValue(ctx, infra.SessionKey, c.Get(infra.SessionKey.String()))
	return ctx
}

// GetAndValidateParam gets the param value and validates it, returning a validation error in case it's invalid
func (r *reqUtils) GetAndValidateParam(c echo.Context, paramName string, errorMessage string) (paramValue string, err error) {
	paramValue = c.Param(paramName)

	if strings.TrimSpace(paramValue) == "" {
		return paramValue, resterrors.NewUnprocessableEntity(errorMessage)
	}

	return paramValue, nil
}

// GetPagingParams gets the standard paging params from the URL, returning how much data to take and skip
func (r *reqUtils) GetPagingParams(c echo.Context, pageParameter, quantityParameter string) (take int64, skip int64) {

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

	return r.GetTakeSkipFromPageQuantity(page, quantity)
}

func (r *reqUtils) GetTakeSkipFromPageQuantity(page, quantity int64) (take, skip int64) {
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
