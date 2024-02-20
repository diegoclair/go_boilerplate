package routeutils

import (
	"context"

	"github.com/labstack/echo/v4"
)

// IRoute interface for register routes
type IRoute interface {
	RegisterRoutes(groups *EchoGroups)
}

// Utils aggregates the request and response utils
type Utils interface {
	Resp() ResponseUtils
	Req() RequestUtils
}

// RequestUtils aggregates the request utils
type RequestUtils interface {
	// GetContext returns a filled ctx with the account uuid and session code if route has access token
	GetContext(c echo.Context) (ctx context.Context)
	GetAndValidateParam(c echo.Context, paramName string, errorMessage string) (paramValue string, err error)
	GetPagingParams(c echo.Context, pageParameter, quantityParameter string) (take int64, skip int64)
	GetTakeSkipFromPageQuantity(page, quantity int64) (take, skip int64)
}

// ResponseUtils aggregates the response utils
type ResponseUtils interface {
	ResponseNoContent(c echo.Context) error
	ResponseCreated(c echo.Context) error
	ResponseAPIOk(c echo.Context, data interface{}) error
	ResponseNotFoundError(c echo.Context, err error) error
	ResponseBadRequestError(c echo.Context, err error) error
	ResponseUnauthorizedError(c echo.Context, err error) error
	ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error
	HandleAPIError(c echo.Context, errorToHandle error) (err error)
}
