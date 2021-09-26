package service

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/util/errors"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/gommon/log"
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

func (s *accountService) CreateAccount(account entity.Account) (err error) {

	log.Info("CreateAccount: Process Started")
	defer log.Info("CreateAccount: Process Finished")

	account.CPF, err = s.svc.cipher.Encrypt(account.CPF)
	if err != nil {
		log.Error("CreateAccount: ", err)
		return err
	}

	_, err = s.svc.dm.Account().GetAccountByDocument(account.CPF)
	if err != nil && !errors.SQLNotFound(err.Error()) {
		log.Error("CreateAccount: ", err)
		return err
	} else if err == nil {
		log.Error("CreateAccount: The document number is already in use")
		return resterrors.NewConflictError("The cpf is already in use")
	}

	account.Secret = s.getHashedPassword(account.Secret)
	account.UUID = uuid.NewV4().String()

	err = s.svc.dm.Account().CreateAccount(account)
	if err != nil {
		log.Error("CreateAccount: ", err)
		return err
	}

	return nil
}

func (s *accountService) getHashedPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *accountService) GetAccounts() (accounts []entity.Account, err error) {

	log.Info("GetAccounts: Process Started")
	defer log.Info("GetAccounts: Process Finished")

	accounts, err = s.svc.dm.Account().GetAccounts()
	if err != nil {
		log.Error("GetAccounts: ", err)
		return accounts, err
	}

	for i := 0; i < len(accounts); i++ {
		_, err = s.svc.cipher.DecryptStruct(&accounts[i])
		if err != nil {
			log.Error("GetAccounts: error to decrypt account struct: ", err)
			return accounts, err
		}
	}

	return accounts, nil
}

func (s *accountService) GetAccountByUUID(accountUUID string) (account entity.Account, err error) {

	log.Info("GetAccountByUUID: Process Started")
	defer log.Info("GetAccountByUUID: Process Finished")

	account, err = s.svc.dm.Account().GetAccountByUUID(accountUUID)
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
