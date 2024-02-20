package transferroute

import (
	"sync"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/go_utils-lib/v2/validator"

	"github.com/labstack/echo/v4"
)

var (
	instance *Handler
	once     sync.Once
)

type Handler struct {
	transferService contract.TransferService
	utils           routeutils.Utils
	structValidator validator.Validator
}

func NewHandler(transferService contract.TransferService, utils routeutils.Utils, structValidator validator.Validator) *Handler {
	once.Do(func() {
		instance = &Handler{
			transferService: transferService,
			utils:           utils,
			structValidator: structValidator,
		}
	})

	return instance
}

func (s *Handler) handleAddTransfer(c echo.Context) error {
	input := viewmodel.TransferReq{}

	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate(s.structValidator)
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

func (s *Handler) handleGetTransfers(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	take, skip := s.utils.Req().GetPagingParams(c, "page", "quantity")

	transfers, totalRecords, err := s.transferService.GetTransfers(ctx, take, skip)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := []viewmodel.TransferResp{}
	for _, transfer := range transfers {
		resp := viewmodel.TransferResp{}
		resp.FillFromEntity(transfer)
		response = append(response, resp)
	}

	responsePaginated := routeutils.BuildPaginatedResult(response, skip, take, totalRecords)

	return s.utils.Resp().ResponseAPIOk(c, responsePaginated)
}
