package transferroute

import (
	"sync"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	transferService contract.TransferService
	mapper          mapper.Mapper
}

func NewController(transferService contract.TransferService, mapper mapper.Mapper) *Controller {
	once.Do(func() {
		instance = &Controller{
			transferService: transferService,
			mapper:          mapper,
		}
	})
	return instance
}

func (s *Controller) handleAddTransfer(c echo.Context) error {

	input := viewmodel.TransferReq{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
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

	ctx := routeutils.GetContext(c)
	transfers, err := s.transferService.GetTransfers(ctx)
	if err != nil {

		return routeutils.HandleAPIError(c, err)
	}

	response := []viewmodel.TransferResp{}
	err = s.mapper.From(transfers).To(&response)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}
	return routeutils.ResponseAPIOK(c, response)
}
