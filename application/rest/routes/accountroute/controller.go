package accountroute

import (
	"sync"

	"github.com/IQ-tech/go-mapper"
	"github.com/diegoclair/go_boilerplate/application/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/application/rest/viewmodel"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/domain/service"

	"github.com/labstack/echo/v4"
)

var (
	instance *Controller
	once     sync.Once
)

type Controller struct {
	accountService service.AccountService
	mapper         mapper.Mapper
}

func NewController(accountService service.AccountService, mapper mapper.Mapper) *Controller {
	once.Do(func() {
		instance = &Controller{
			accountService: accountService,
			mapper:         mapper,
		}
	})
	return instance
}

func (s *Controller) handleAddAccount(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	input := viewmodel.AddAccount{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	account := entity.Account{
		Name:   input.Name,
		CPF:    input.CPF,
		Secret: input.Secret,
	}

	err = s.accountService.CreateAccount(ctx, account)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}
	return routeutils.ResponseCreated(c)
}

func (s *Controller) handleAddBalance(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	input := viewmodel.AddBalance{}
	err := c.Bind(&input)
	if err != nil {
		return routeutils.ResponseBadRequestError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	accountUUID, err := routeutils.GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	err = s.accountService.AddBalance(ctx, accountUUID, input.Amount)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}
	return routeutils.ResponseCreated(c)
}

func (s *Controller) handleGetAccounts(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	take, skip := routeutils.GetPagingParams(c, "page", "quantity")

	accounts, totalRecords, err := s.accountService.GetAccounts(ctx, take, skip)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := []viewmodel.Account{}
	err = s.mapper.From(accounts).To(&response)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	responsePaginated := routeutils.BuildPaginatedResult(response, skip, take, totalRecords)

	return routeutils.ResponseAPIOK(c, responsePaginated)
}

func (s *Controller) handleGetAccountByID(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	accountUUID, err := routeutils.GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.Account{}
	err = s.mapper.From(account).To(&response)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	return routeutils.ResponseAPIOK(c, response)
}

func (s *Controller) handleGetAccountBalanceByID(c echo.Context) error {

	ctx := routeutils.GetContext(c)

	accountUUID, err := routeutils.GetAndValidateParam(c, "account_id", "Invalid account_id")
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return routeutils.HandleAPIError(c, err)
	}

	response := viewmodel.Account{
		Balance: account.Balance,
	}

	return routeutils.ResponseAPIOK(c, response)
}
