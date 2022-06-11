package routeutils

import (
	"math"
	"net/http"

	"github.com/diegoclair/go-boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// BuildPaginatedResult returns a paginatedResult instance
func BuildPaginatedResult(list interface{}, skip int64, take int64, totalRecords int64) viewmodel.PaginatedResult {
	return viewmodel.PaginatedResult{
		List: list,
		Pagination: viewmodel.ReturnPagination{
			CurrentPage:    (skip / take) + 1,
			RecordsPerPage: take,
			TotalRecords:   totalRecords,
			TotalPages:     int64(math.Ceil(float64(totalRecords) / float64(take))),
		},
	}
}

const ErrorMessageServiceUnavailable = "Service temporarily unavailable"

func ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func ResponseCreated(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func ResponseAPIOK(c echo.Context, data interface{}) error {
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
