package accountroute

import (
	"sync"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"

	echo "github.com/labstack/echo/v4"
)

var (
	instance *Handler
	Once     sync.Once
)

type Handler struct {
	accountService contract.AccountService
}

func NewHandler(accountService contract.AccountService) *Handler {
	Once.Do(func() {
		instance = &Handler{
			accountService: accountService,
		}
	})

	return instance
}

func (s *Handler) handleAddAccount(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.AddAccount{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	err = s.accountService.CreateAccount(ctx, input.ToDto())
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	return routeutils.ResponseCreated(c)
}

func (s *Handler) handleAddBalance(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	input := viewmodel.AddBalance{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	accountUUID, err := routeutils.GetAndValidateParam(c, "account_uuid", "account_uuid is required")
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	err = s.accountService.AddBalance(ctx, input.ToDto(accountUUID))
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	return routeutils.ResponseCreated(c)
}

func (s *Handler) handleGetAccounts(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	take, skip := routeutils.GetPagingParams(c, "page", "quantity")

	accounts, totalRecords, err := s.accountService.GetAccounts(ctx, take, skip)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := []viewmodel.AccountResponse{}
	for _, account := range accounts {
		item := viewmodel.AccountResponse{}
		item.FillFromEntity(account)
		response = append(response, item)
	}

	responsePaginated := viewmodel.BuildPaginatedResponse(response, skip, take, totalRecords)

	return routeutils.ResponseAPIOk(c, responsePaginated)
}

func (s *Handler) handleGetAccountByID(c echo.Context) error {
	ctx := routeutils.GetContext(c)

	accountUUID, err := routeutils.GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.AccountResponse{}
	response.FillFromEntity(account)

	return routeutils.ResponseAPIOk(c, response)
}
