package service

import (
	"context"
	"fmt"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/crypto"
	utilerrors "github.com/diegoclair/go-boilerplate/util/errors"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/twinj/uuid"
)

type accountService struct {
	svc *Service
}

func newAccountService(svc *Service) AccountService {
	return &accountService{
		svc: svc,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, account entity.Account) (err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	_, err = s.svc.dm.Account().GetAccountByDocument(ctx, account.CPF)
	if err != nil && !utilerrors.SQLNotFound(err.Error()) {
		log.Error(err)
		return err
	} else if err == nil {
		log.Error("The document number is already in use")
		return resterrors.NewConflictError("The cpf is already in use")
	}

	account.Secret, err = crypto.HashPassword(account.Secret)
	if err != nil {
		log.Error(err)
		return err
	}
	account.UUID = uuid.NewV4().String()

	err = s.svc.dm.Account().CreateAccount(ctx, account)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//TODO: fix floating point
func (s *accountService) AddBalance(ctx context.Context, accountUUID string, amount float64) (err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		log.Error("error to get account", err)
		return err
	}
	balance := account.Balance + amount

	err = s.svc.dm.Account().UpdateAccountBalance(ctx, account.ID, balance)
	if err != nil {
		log.Error("error to update account balance", err)
		return err
	}

	return nil
}

func (s *accountService) GetAccounts(ctx context.Context, take, skip int64) (accounts []entity.Account, totalRecords int64, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	accounts, totalRecords, err = s.svc.dm.Account().GetAccounts(ctx, take, skip)
	if err != nil {
		log.Error(err)
		return accounts, totalRecords, err
	}

	log.Info(fmt.Sprintf("Found %d accounts", totalRecords))

	return accounts, totalRecords, nil
}

func (s *accountService) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started with accountUUID: ", accountUUID)
	defer log.Info("Process Finished for accountUUID: ", accountUUID)

	account, err = s.svc.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		log.Error(err)
		return account, err
	}

	return account, nil
}
