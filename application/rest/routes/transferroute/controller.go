package transferroute

import (
	"sync"

	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/contract"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	transferService contract.TransferService
	utils           routeutils.Utils
}

func NewController(transferService contract.TransferService, utils routeutils.Utils) *Controller {
	once.Do(func() {
		instance = &Controller{
			transferService: transferService,
			utils:           utils,
		}
	})
	return instance
}

func (s *Controller) handleAddTransfer(c echo.Context) error {
	input := viewmodel.TransferReq{}

	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	appContext := s.utils.Req().GetContext(c)

	err = s.transferService.CreateTransfer(appContext, input.ToEntity())
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	return s.utils.Resp().ResponseCreated(c)
}

func (s *Controller) handleGetTransfers(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	transfers, err := s.transferService.GetTransfers(ctx)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := []viewmodel.TransferResp{}
	for _, transfer := range transfers {
		resp := viewmodel.TransferResp{}
		resp.FillFromEntity(transfer)
		response = append(response, resp)
	}

	return s.utils.Resp().ResponseAPIOK(c, response)
}
