package authroute

import (
	"sync"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go-boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go-boilerplate/domain/service"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	authService service.AuthService
	mapper      mapper.Mapper
}

func NewController(authService service.AuthService, mapper mapper.Mapper) *Controller {
	once.Do(func() {
		instance = &Controller{
			authService: authService,
			mapper:      mapper,
		}
	})
	return instance
}

func (s *Controller) handleLogin(c echo.Context) error {

	input := viewmodel.Login{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}
	err = input.Validate()
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	auth, err := s.authService.Login(input.CPF, input.Secret)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.AuthResponse{}
	err = s.mapper.From(auth).To(&response)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	return routeutils.ResponseAPIOK(c, response)
}
