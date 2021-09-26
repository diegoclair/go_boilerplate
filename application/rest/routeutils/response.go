package routeutils

import (
	"net/http"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const ErrorMessageServiceUnavailable = "Service temporarily unavailable"

// ResponseNoContent returns a standard API success with no content response
func ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// ResponseCreated returns a standard API successful as a result with created code
func ResponseCreated(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

// ResponseAPIOK returns a standard API success response
func ResponseAPIOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// ResponseAPINotFoundError returns a standard API not found error
func ResponseAPINotFoundError(c echo.Context) error {
	return ResponseAPIError(c, http.StatusNotFound, "Not Found", "", nil)
}

// ResponseAPIError returns a standard API error
func ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error {
	returnValue := resterrors.NewRestError(message, status, err, causes)
	return c.JSON(status, returnValue)
}

func HandleAPIError(c echo.Context, errorToHandle error) (err error) {
	statusCode := http.StatusServiceUnavailable
	errorMessage := ErrorMessageServiceUnavailable

	if errorToHandle != nil {
		log.Error("HandleAPIError: ", errorToHandle)

		errorString := errorToHandle.Error()

		restErr, ok := errorToHandle.(resterrors.RestErr)
		if !ok {
			return ResponseAPIError(c, statusCode, errorMessage, errorString, nil)
		}

		return c.JSON(restErr.StatusCode(), restErr)

	}

	return ResponseAPIError(c, statusCode, errorMessage, "", nil)
}
