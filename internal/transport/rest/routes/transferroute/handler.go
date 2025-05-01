package transferroute

import (
	"sync"

	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/internal/transport/rest/viewmodel"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	transferService contract.TransferApp
}

func NewHandler(transferService contract.TransferApp) *Handler {
	Once.Do(func() {
		instance = &Handler{
			transferService: transferService,
		}
	})

	return instance
}

func (s *Handler) handleAddTransfer(c echo.Context) error {
	input := viewmodel.TransferReq{}

	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseInvalidRequestBody(c, err)
	}

	appContext := routeutils.GetContext(c)

	err = s.transferService.CreateTransfer(appContext, input.ToDto())
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	return routeutils.ResponseCreated(c)
}

func (s *Handler) handleGetTransfers(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	take, skip := routeutils.GetPagingParams(c, "page", "quantity")

	transfers, totalRecords, err := s.transferService.GetTransfers(ctx, take, skip)
	if err != nil {
		return routeutils.HandleError(c, err)
	}

	response := []viewmodel.TransferResp{}
	for _, transfer := range transfers {
		resp := viewmodel.TransferResp{}
		resp.FillFromEntity(transfer)
		response = append(response, resp)
	}

	responsePaginated := viewmodel.BuildPaginatedResponse(response, skip, take, totalRecords)

	return routeutils.ResponseAPIOk(c, responsePaginated)
}
