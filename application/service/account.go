package service

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/twinj/uuid"
)

type accountService struct {
	svc *service
}

func newAccountService(svc *service) contract.AccountService {
	return &accountService{
		svc: svc,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, input dto.AccountInput) (err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	account, err := input.ToEntityValidate(ctx, s.svc.validator)
	if err != nil {
		s.svc.log.Errorf(ctx, "error or invalid input: %s", err.Error())
		return err
	}

	_, err = s.svc.dm.Account().GetAccountByDocument(ctx, account.CPF)
	if err != nil && !mysqlutils.SQLNotFound(err.Error()) {
		s.svc.log.Errorf(ctx, "error to get account by document: %s", err.Error())
		return err
	} else if err == nil {
		s.svc.log.Error(ctx, "The document number is already in use")
		return resterrors.NewConflictError("The cpf is already in use")
	}

	account.Password, err = s.svc.crypto.HashPassword(account.Password)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to hash password: %s", err.Error())
		return err
	}
	account.UUID = uuid.NewV4().String()

	_, err = s.svc.dm.Account().CreateAccount(ctx, account)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to create account: %s", err.Error())
		return err
	}

	return nil
}

func (s *accountService) AddBalance(ctx context.Context, input dto.AddBalanceInput) (err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	err = input.Validate(ctx, s.svc.validator)
	if err != nil {
		s.svc.log.Errorf(ctx, "error or invalid input: %s", err.Error())
		return err
	}

	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, input.AccountUUID)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to get account by uuid: %s", err.Error())
		return err
	}

	account.AddBalance(input.Amount)

	err = s.svc.dm.Account().UpdateAccountBalance(ctx, account.ID, account.Balance)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to update account balance: %s", err.Error())
		return err
	}

	return nil
}

func (s *accountService) GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	accounts, totalRecords, err = s.svc.dm.Account().GetAccounts(ctx, take, skip)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to get accounts: %s", err.Error())
		return accounts, totalRecords, err
	}

	s.svc.log.Infof(ctx, "Found %d accounts", totalRecords)

	return accounts, totalRecords, nil
}

func (s *accountService) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {
	s.svc.log.Infof(ctx, "Process Started with accountUUID: %s", accountUUID)
	defer s.svc.log.Infof(ctx, "Process Finished for accountUUID: %s", accountUUID)

	account, err = s.svc.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		s.svc.log.Errorf(ctx, "error to get account by uuid: %s", err.Error())
		return account, err
	}

	return account, nil
}
