package routeutils

import "github.com/diegoclair/goswag/models"

// EchoGroups is the struct that holds the echo groups for the routes
type EchoGroups struct {
	// AppGroup is the group for public routes
	AppGroup models.EchoGroup
	// PrivateGroup is the group for routes that need to be authenticated (login required)
	PrivateGroup models.EchoGroup
}
