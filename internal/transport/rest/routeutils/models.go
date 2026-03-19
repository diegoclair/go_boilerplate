package routeutils

import (
	"net/http"

	"github.com/diegoclair/apperr/httpmap"
	"github.com/diegoclair/goswag/models"
)

// EchoGroups is the struct that holds the echo groups for the routes
type EchoGroups struct {
	// AppGroup is the group for public routes
	AppGroup models.EchoGroup
	// PrivateGroup is the group for routes that need to be authenticated (login required)
	PrivateGroup models.EchoGroup
}

// DefaultSwaggerErrors returns the standard error responses for Swagger documentation.
func DefaultSwaggerErrors() []models.ReturnType {
	return []models.ReturnType{
		{StatusCode: http.StatusBadRequest, Body: httpmap.ErrorResponse{}},
		{StatusCode: http.StatusUnauthorized, Body: httpmap.ErrorResponse{}},
		{StatusCode: http.StatusForbidden, Body: httpmap.ErrorResponse{}},
		{StatusCode: http.StatusNotFound, Body: httpmap.ErrorResponse{}},
		{StatusCode: http.StatusConflict, Body: httpmap.ErrorResponse{}},
		{StatusCode: http.StatusInternalServerError, Body: httpmap.ErrorResponse{}},
	}
}
