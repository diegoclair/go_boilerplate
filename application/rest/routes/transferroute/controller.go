package transferroute

import (
	"sync"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go-boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go-boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/domain/service"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	transferService service.TransferService
	mapper          mapper.Mapper
}

func NewController(transferService service.TransferService, mapper mapper.Mapper) *Controller {
	once.Do(func() {
		instance = &Controller{
			transferService: transferService,
			mapper:          mapper,
		}
	})
	return instance
}

func (s *Controller) handleAddTransfer(c echo.Context) error {

	input := viewmodel.Transfer{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	transfer := entity.Transfer{}
	err = s.mapper.From(input).To(&transfer)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	appContext := routeutils.GetContext(c)
	err = s.transferService.CreateTransfer(appContext, transfer)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}
	return routeutils.ResponseCreated(c)
}

func (s *Controller) handleGetTransfers(c echo.Context) error {

	appContext := routeutils.GetContext(c)
	transfers, err := s.transferService.GetTransfers(appContext)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := []viewmodel.Transfer{}
	err = s.mapper.From(transfers).To(&response)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	return routeutils.ResponseAPIOK(c, response)
}
