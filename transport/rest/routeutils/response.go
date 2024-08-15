package routeutils

import (
	"net/http"

	"github.com/diegoclair/go_utils/resterrors"
	echo "github.com/labstack/echo/v4"
)

const ErrorMessageServiceUnavailable = "Service temporarily unavailable"

func ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func ResponseCreated(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func ResponseAPIOk(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func ResponseNotFoundError(c echo.Context, err error) error {
	return ResponseAPIError(c, http.StatusNotFound, "Not Found", err.Error(), nil)
}

func ResponseBadRequestError(c echo.Context, err error) error {
	return ResponseAPIError(c, http.StatusBadRequest, "Bad request", err.Error(), nil)
}

func ResponseUnauthorizedError(c echo.Context, err error) error {
	return ResponseAPIError(c, http.StatusUnauthorized, "Unauthorized", err.Error(), nil)
}

func ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error {
	returnValue := resterrors.NewRestError(message, status, err, causes)
	return c.JSON(status, returnValue)
}

func HandleAPIError(c echo.Context, errorToHandle error) (err error) {
	statusCode := http.StatusServiceUnavailable
	errorMessage := ErrorMessageServiceUnavailable

	if errorToHandle != nil {
		errorString := errorToHandle.Error()

		restErr, ok := errorToHandle.(resterrors.RestErr)
		if !ok {
			return ResponseAPIError(c, statusCode, errorMessage, errorString, nil)
		}

		return c.JSON(restErr.StatusCode(), restErr)

	}

	return ResponseAPIError(c, statusCode, errorMessage, "", nil)
}
