package accountroute

import (
	"sync"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/domain/account"
	"github.com/diegoclair/go_boilerplate/transport/rest/routeutils"
	"github.com/diegoclair/go_boilerplate/transport/rest/viewmodel"
	"github.com/diegoclair/go_utils/validator"

	"github.com/labstack/echo/v4"
)

var (
	instance *Handler
	once     sync.Once
)

type Handler struct {
	accountService contract.AccountService
	utils          routeutils.Utils
	validator      validator.Validator
}

func NewHandler(accountService contract.AccountService, utils routeutils.Utils, validator validator.Validator) *Handler {
	once.Do(func() {
		instance = &Handler{
			accountService: accountService,
			utils:          utils,
			validator:      validator,
		}
	})

	return instance
}

func (s *Handler) handleAddAccount(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.AddAccount{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate(s.validator)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	account := account.Account{
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

func (s *Handler) handleAddBalance(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	input := viewmodel.AddBalance{}
	err := c.Bind(&input)
	if err != nil {
		return s.utils.Resp().ResponseBadRequestError(c, err)
	}

	err = input.Validate(s.validator)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	accountUUID, err := s.utils.Req().GetAndValidateParam(c, "account_uuid", "account_uuid is required")
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	err = s.accountService.AddBalance(ctx, accountUUID, input.Amount)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	return s.utils.Resp().ResponseCreated(c)
}

func (s *Handler) handleGetAccounts(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	take, skip := s.utils.Req().GetPagingParams(c, "page", "quantity")

	accounts, totalRecords, err := s.accountService.GetAccounts(ctx, take, skip)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := []viewmodel.AccountResponse{}
	for _, account := range accounts {
		item := viewmodel.AccountResponse{}
		item.FillFromEntity(account)
		response = append(response, item)
	}

	responsePaginated := routeutils.BuildPaginatedResult(response, skip, take, totalRecords)

	return s.utils.Resp().ResponseAPIOk(c, responsePaginated)
}

func (s *Handler) handleGetAccountByID(c echo.Context) error {
	ctx := s.utils.Req().GetContext(c)

	accountUUID, err := s.utils.Req().GetAndValidateParam(c, "account_uuid", "Invalid account_uuid")
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	account, err := s.accountService.GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		return s.utils.Resp().HandleAPIError(c, err)
	}

	response := viewmodel.AccountResponse{}
	response.FillFromEntity(account)

	return s.utils.Resp().ResponseAPIOk(c, response)
}
