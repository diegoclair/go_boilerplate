package service

import (
	"context"
	"fmt"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/crypto"
	"github.com/diegoclair/go-boilerplate/util/errors"
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
	log.Info("CreateAccount: Process Started")
	defer log.Info("CreateAccount: Process Finished")

	account.CPF, err = s.svc.cipher.Encrypt(account.CPF)
	if err != nil {
		log.Error("CreateAccount: ", err)
		return err
	}

	_, err = s.svc.dm.Account().GetAccountByDocument(ctx, account.CPF)
	if err != nil && !errors.SQLNotFound(err.Error()) {
		log.Error("CreateAccount: ", err)
		return err
	} else if err == nil {
		log.Error("CreateAccount: The document number is already in use")
		return resterrors.NewConflictError("The cpf is already in use")
	}

	account.Secret = crypto.GetMd5(account.Secret)
	account.UUID = uuid.NewV4().String()

	err = s.svc.dm.Account().CreateAccount(ctx, account)
	if err != nil {
		log.Error("CreateAccount: ", err)
		return err
	}

	return nil
}

//TODO: search about floating point
func (s *accountService) AddBalance(ctx context.Context, accountUUID string, amount float64) (err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("AddBalance: Process Started")
	defer log.Info("AddBalance: Process Finished")

	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		log.Error("AddBalance: error to get account", err)
		return err
	}
	account.Balance += amount

	err = s.svc.dm.Account().UpdateAccountBalance(ctx, account)
	if err != nil {
		log.Error("AddBalance: error to update account balance", err)
		return err
	}

	return nil
}

func (s *accountService) GetAccounts(ctx context.Context) (accounts []entity.Account, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("GetAccounts: Process Started")
	defer log.Info("GetAccounts: Process Finished")

	accounts, err = s.svc.dm.Account().GetAccounts(ctx)
	if err != nil {
		log.Error("GetAccounts: ", err)
		return accounts, err
	}

	for i := 0; i < len(accounts); i++ {
		fmt.Println(accounts[i])
		// _, err = s.svc.cipher.DecryptStruct(&accounts[i])
		// if err != nil {
		// 	log.Error("GetAccounts: error to decrypt account struct: ", err)
		// 	return accounts, err
		// }
	}

	return accounts, nil
}

func (s *accountService) GetAccountByUUID(ctx context.Context, accountUUID string) (account entity.Account, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("GetAccountByUUID: Process Started")
	defer log.Info("GetAccountByUUID: Process Finished")

	account, err = s.svc.dm.Account().GetAccountByUUID(ctx, accountUUID)
	if err != nil {
		log.Error("GetAccountByUUID: ", err)
		return account, err
	}

	_, err = s.svc.cipher.DecryptStruct(&account)
	if err != nil {
		log.Error("GetAccountByUUID: error to decrypt account struct: ", err)
		return account, err
	}

	return account, nil
}
