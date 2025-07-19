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

// ArrayConverter defines a function that converts a string to a specific type
type ArrayConverter[T any] func(value string) (T, error)

// GetRequiredParam gets a parameter value, converts it using the provided converter, and validates it's not zero value
func GetRequiredParam[T comparable](rawValue string, converter ArrayConverter[T], errorMessage string) (T, error) {
	var zero T

	// Check for empty string first
	if strings.TrimSpace(rawValue) == "" {
		return zero, resterrors.NewUnprocessableEntity(errorMessage)
	}

	// Convert the value using the same converter as arrays
	value, err := converter(rawValue)
	if err != nil {
		return zero, resterrors.NewUnprocessableEntity(errorMessage)
	}

	// Check if result is zero value (Go's zero value check)
	if value == zero {
		return zero, resterrors.NewUnprocessableEntity(errorMessage)
	}

	return value, nil
}

// Convenience functions using the generic base function with existing converters
func GetRequiredInt64PathParam(c echo.Context, paramName string, errorMessage string) (int64, error) {
	return GetRequiredParam(c.Param(paramName), Int64Converter, errorMessage)
}

func GetRequiredStringPathParam(c echo.Context, paramName string, errorMessage string) (string, error) {
	return GetRequiredParam(c.Param(paramName), StringConverter, errorMessage)
}

func GetRequiredStringQueryParam(c echo.Context, paramName string, errorMessage string) (string, error) {
	return GetRequiredParam(c.QueryParam(paramName), StringConverter, errorMessage)
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

// GetArrayParam gets an array of any type from a string using a converter function and separator
func GetArrayParam[T any](rawValue, separator string, converter ArrayConverter[T]) ([]T, error) {
	if strings.TrimSpace(rawValue) == "" {
		return []T{}, nil
	}

	// Split and convert each item using the provided converter
	items := strings.Split(rawValue, separator)
	result := make([]T, 0, len(items))

	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			value, err := converter(trimmed)
			if err != nil {
				return nil, resterrors.NewUnprocessableEntity("Invalid value: " + trimmed + " - " + err.Error())
			}
			result = append(result, value)
		}
	}

	return result, nil
}

// String converter - no conversion needed
func StringConverter(value string) (string, error) {
	return value, nil
}

// Int64 converter
func Int64Converter(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// Int converter
func IntConverter(value string) (int, error) {
	return strconv.Atoi(value)
}

// Convenience functions using the generic base function
func GetStringArrayQueryParam(c echo.Context, paramName, separator string) []string {
	result, _ := GetArrayParam(c.QueryParam(paramName), separator, StringConverter)
	return result
}

func GetInt64ArrayQueryParam(c echo.Context, paramName, separator string) ([]int64, error) {
	return GetArrayParam(c.QueryParam(paramName), separator, Int64Converter)
}

func GetIntArrayQueryParam(c echo.Context, paramName, separator string) ([]int, error) {
	return GetArrayParam(c.QueryParam(paramName), separator, IntConverter)
}
