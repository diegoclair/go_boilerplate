package accountroute

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
	accountService contract.AccountService
	mapper         mapper.Mapper
	utils          routeutils.Utils
}

func NewController(accountService contract.AccountService, mapper mapper.Mapper, utils routeutils.Utils) *Controller {
	once.Do(func() {
		instance = &Controller{
			accountService: accountService,
			mapper:         mapper,
			utils:          utils,
		}
	})
	return instance
}

func (s *Controller) handleAddAccount(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.AddAccount{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	account := entity.Account{
		Name:     input.Name,
		CPF:      input.CPF,
		Password: input.Password,
	}

	err = s.accountService.CreateAccount(ctx, account)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}
	return s.utils.Resp().ResponseCreated(c)
}

func (s *Controller) handleAddBalance(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.AddBalance{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate()
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	accountUUID, err := s.utils.Req().GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	err = s.accountService.AddBalance(ctx, accountUUID, input.Amount)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}
	return s.utils.Resp().ResponseCreated(c)
}

func (s *Controller) handleGetAccounts(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	take, skip := s.utils.Req().GetPagingParams(c, "page", "quantity")

	accounts, totalRecords, err := s.accountService.GetAccounts(ctx, take, skip)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := []viewmodel.Account{}
	err = s.mapper.From(accounts).To(&response)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	responsePaginated := s.utils.Resp().BuildPaginatedResult(response, skip, take, totalRecords)

	return s.utils.Resp().ResponseAPIOK(c, responsePaginated)
}

func (s *Controller) handleGetAccountByID(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	accountUUID, err := s.utils.Req().GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := viewmodel.Account{}
	err = s.mapper.From(account).To(&response)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	return s.utils.Resp().ResponseAPIOK(c, response)
}

func (s *Controller) handleGetAccountBalanceByID(c echo.Context) error {

	ctx := s.utils.Req().GetContext(c)

	accountUUID, err := s.utils.Req().GetAndValidateParam(c, "account_id", "Invalid account_id")
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := viewmodel.Account{
		Balance: account.Balance,
	}

	return s.utils.Resp().ResponseAPIOK(c, response)
}
