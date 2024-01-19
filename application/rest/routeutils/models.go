package routeutils

import "github.com/labstack/echo/v4"

// EchoGroups is the struct that holds the echo groups for the routes
type EchoGroups struct {
	// AppGroup is the group for public routes
	AppGroup *echo.Group
	// PrivateGroup is the group for routes that need to be authenticated (login required)
	PrivateGroup *echo.Group
	// AdminGroup is the group for routes that need to be authenticated and have admin role (login required)
	AdminGroup *echo.Group
	// SuperAdminGroup is the group for routes that need to be authenticated and have super admin role (login required)
	SuperAdminGroup *echo.Group
}
