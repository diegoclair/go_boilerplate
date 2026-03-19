package routeutils

import (
	"net/http"

	"github.com/diegoclair/apperr/httpmap"
	echo "github.com/labstack/echo/v4"
)

const ErrorMessageServiceUnavailable = "Service temporarily unavailable"

func ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// ResponseCreated returns a 201 Created response and a json body if data is provided
func ResponseCreated(c echo.Context, data ...interface{}) error {
	if len(data) > 0 {
		return c.JSON(http.StatusCreated, data[0])
	}

	return c.NoContent(http.StatusCreated)
}

func ResponseAPIOk(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func ResponseInvalidRequestBody(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, httpmap.ErrorResponse{
		Message:    "invalid request body",
		StatusCode: http.StatusBadRequest,
		Error:      http.StatusText(http.StatusBadRequest),
	})
}

func HandleError(c echo.Context, errorToHandle error) error {
	status, body := httpmap.ToHTTP(errorToHandle)
	return c.JSON(status, body)
}
