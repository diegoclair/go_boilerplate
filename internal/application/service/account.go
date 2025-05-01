package service

import (
	"context"
	"errors"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/twinj/uuid"
)

type accountService struct {
	crypto    contract.Crypto
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
}

func newAccountService(infra domain.Infrastructure) *accountService {
	return &accountService{
		crypto:    infra.Crypto(),
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
	}
}

func (s *accountService) CreateAccount(ctx context.Context, input dto.AccountInput) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	account, err := input.ToEntityValidate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	_, err = s.dm.Account().GetAccountByDocument(ctx, account.CPF)
	if err != nil && !mysqlutils.SQLNotFound(err.Error()) {
		s.log.Errorw(ctx, "error to get account by document", logger.Err(err))
		return err
	} else if err == nil {
		s.log.Error(ctx, "The document number is already in use")
		return resterrors.NewConflictError("The cpf is already in use")
	}

	account.Password, err = s.crypto.HashPassword(account.Password)
	if err != nil {
		s.log.Errorw(ctx, "error to hash password", logger.Err(err))
		return err
	}
	account.UUID = uuid.NewV4().String()

	_, err = s.dm.Account().CreateAccount(ctx, account)
	if err != nil {
		s.log.Errorw(ctx, "error to create account", logger.Err(err))
		return err
	}

	return nil
}

func (s *accountService) AddBalance(ctx context.Context, input dto.AddBalanceInput) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	account, err := s.dm.Account().GetAccountByUUID(ctx, input.AccountUUID)
	if err != nil {
		s.log.Errorw(ctx, "error to get account by uuid", logger.Err(err))
		return err
	}

	account.AddBalance(input.Amount)

	err = s.dm.Account().UpdateAccountBalance(ctx, account.ID, account.Balance)
	if err != nil {
		s.log.Errorw(ctx, "error to update account balance", logger.Err(err))
		return err
	}

	return nil
}

func (s *accountService) GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	accounts, totalRecords, err = s.dm.Account().GetAccounts(ctx, take, skip)
	if err != nil {
		s.log.Errorw(ctx, "error to get accounts", logger.Err(err))
		return accounts, totalRecords, err
	}

	s.log.Infof(ctx, "Found %d accounts", totalRecords)

	return accounts, totalRecords, nil
}

func (s *accountService) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {
	s.log.Infof(ctx, "Process Started with accountUUID: %s", accountUUID)
	defer s.log.Infof(ctx, "Process Finished for accountUUID: %s", accountUUID)

	account, err = s.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		s.log.Errorw(ctx, "error to get account by uuid", logger.Err(err))
		return account, err
	}

	return account, nil
}

func (s *accountService) getLoggedAccountUUID(ctx context.Context) (accountUUID string, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	loggedAccountUUID, ok := ctx.Value(infra.AccountUUIDKey).(string)
	if !ok {
		errMsg := "accountUUID should not be empty"
		s.log.Error(ctx, errMsg)
		return accountUUID, errors.New(errMsg)
	}

	return loggedAccountUUID, nil
}

func (s *accountService) GetLoggedAccountID(ctx context.Context) (accountID int64, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	loggedAccountUUID, err := s.getLoggedAccountUUID(ctx)
	if err != nil {
		return accountID, err
	}

	accountID, err = s.dm.Account().GetAccountIDByUUID(ctx, loggedAccountUUID)
	if err != nil {
		s.log.Errorw(ctx, "error to get logged account ID by uuid", logger.Err(err))
		return accountID, err
	}

	return accountID, nil
}

func (s *accountService) GetLoggedAccount(ctx context.Context) (account entity.Account, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	loggedAccountUUID, err := s.getLoggedAccountUUID(ctx)
	if err != nil {
		return account, err
	}

	account, err = s.dm.Account().GetAccountByUUID(ctx, loggedAccountUUID)
	if err != nil {
		s.log.Errorw(ctx, "error to get logged account by uuid", logger.Err(err))
		return account, err
	}

	return account, nil
}
