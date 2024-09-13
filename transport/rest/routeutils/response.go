package routeutils

import (
	"fmt"
	"net/http"

	"github.com/diegoclair/go_utils/resterrors"
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

func ResponseUnauthorizedError(c echo.Context, errMsg string) error {
	err := resterrors.NewUnauthorizedError(errMsg)
	return c.JSON(err.StatusCode(), err)
}

func ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error {
	returnValue := resterrors.NewRestError(message, status, err, causes)
	return c.JSON(status, returnValue)
}

func ResponseInvalidRequestBody(c echo.Context, err error) error {
	e := resterrors.NewBadRequestError("Invalid request body", err)
	return c.JSON(e.StatusCode(), e)
}

func HandleError(c echo.Context, errorToHandle error) (err error) {
	statusCode := http.StatusServiceUnavailable
	errorMessage := ErrorMessageServiceUnavailable

	if errorToHandle != nil {
		errorString := errorToHandle.Error()

		restErr, ok := errorToHandle.(resterrors.RestErr)
		if !ok {
			fmt.Println("errorToHandle", errorToHandle, ok)
			return ResponseAPIError(c, statusCode, errorMessage, errorString, nil)
		}

		return c.JSON(restErr.StatusCode(), restErr)

	}

	return ResponseAPIError(c, statusCode, errorMessage, "", nil)
}
