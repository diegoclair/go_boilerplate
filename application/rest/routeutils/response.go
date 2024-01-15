package routeutils

import (
	"math"
	"net/http"

	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/infra/logger"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/echo/v4"
)

type ResponseUtils interface {
	BuildPaginatedResult(list interface{}, skip int64, take int64, totalRecords int64) viewmodel.PaginatedResult
	ResponseNoContent(c echo.Context) error
	ResponseCreated(c echo.Context) error
	ResponseAPIOk(c echo.Context, data interface{}) error
	ResponseNotFoundError(c echo.Context, err error) error
	ResponseBadRequestError(c echo.Context, err error) error
	ResponseUnauthorizedError(c echo.Context, err error) error
	ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error
	HandleAPIError(c echo.Context, errorToHandle error) (err error)
}

type respUtils struct {
	log logger.Logger
}

func newResponseUtils(log logger.Logger) ResponseUtils {
	return &respUtils{log: log}
}

// BuildPaginatedResult returns a paginatedResult instance
func (r *respUtils) BuildPaginatedResult(list interface{}, skip int64, take int64, totalRecords int64) viewmodel.PaginatedResult {
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

func (r *respUtils) ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (r *respUtils) ResponseCreated(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func (r *respUtils) ResponseAPIOk(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func (r *respUtils) ResponseNotFoundError(c echo.Context, err error) error {
	return r.ResponseAPIError(c, http.StatusNotFound, "Not Found", err.Error(), nil)
}

func (r *respUtils) ResponseBadRequestError(c echo.Context, err error) error {
	return r.ResponseAPIError(c, http.StatusBadRequest, "Bad request", err.Error(), nil)
}

func (r *respUtils) ResponseUnauthorizedError(c echo.Context, err error) error {
	return r.ResponseAPIError(c, http.StatusUnauthorized, "Unauthorized", err.Error(), nil)
}

func (r *respUtils) ResponseAPIError(c echo.Context, status int, message string, err string, causes interface{}) error {
	returnValue := resterrors.NewRestError(message, status, err, causes)
	return c.JSON(status, returnValue)
}

func (r *respUtils) HandleAPIError(c echo.Context, errorToHandle error) (err error) {
	statusCode := http.StatusServiceUnavailable
	errorMessage := ErrorMessageServiceUnavailable

	if errorToHandle != nil {
		r.log.Errorf(c.Request().Context(), "HandleAPIError: %s", errorToHandle)

		errorString := errorToHandle.Error()

		restErr, ok := errorToHandle.(resterrors.RestErr)
		if !ok {
			return r.ResponseAPIError(c, statusCode, errorMessage, errorString, nil)
		}

		return c.JSON(restErr.StatusCode(), restErr)

	}

	return r.ResponseAPIError(c, statusCode, errorMessage, "", nil)
}
