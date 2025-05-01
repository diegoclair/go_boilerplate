package routeutils

// IRoute interface for register routes
type IRoute interface {
	RegisterRoutes(groups *EchoGroups)
}
