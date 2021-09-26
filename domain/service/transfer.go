package service

import (
	"context"

	"github.com/diegoclair/go-boilerplate/domain/entity"
	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/gommon/log"
	"github.com/twinj/uuid"
)

type transferService struct {
	svc *Service
}

func newTransferService(svc *Service) TransferService {
	return &transferService{
		svc: svc,
	}
}

func (s *transferService) CreateTransfer(appContext context.Context, transfer entity.Transfer) (err error) {

	log.Info("CreateTransfer: Process Started")
	defer log.Info("CreateTransfer: Process Finished")

	loggedAccountUUID := appContext.Value(auth.AccountUUIDKey)
	account, err := s.svc.dm.MySQL().Account().GetAccountByUUID(loggedAccountUUID.(string))
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	if account.Balance < transfer.Amount {
		return resterrors.NewConflictError("Your account don't have sufficient funds to do this operation")
	}

	destAccount, err := s.svc.dm.MySQL().Account().GetAccountByUUID(transfer.AccountDestinationUUID)
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	transfer.AccountDestinationID = destAccount.ID
	transfer.AccountOriginID = account.ID
	transfer.TransferUUID = uuid.NewV4().String()

	tx, err := s.svc.dm.MySQL().Begin()
	if err != nil {
		log.Error("CreateTransfer: error to get db transaction", err)
		return err
	}
	defer tx.Rollback()

	err = tx.Account().AddTransfer(transfer)
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	account.Balance -= transfer.Amount
	err = tx.Account().UpdateAccountBalance(account)
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	destAccount.Balance += transfer.Amount
	err = tx.Account().UpdateAccountBalance(destAccount)
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error("CreateTransfer: ", err)
		return err
	}

	return nil
}

func (s *transferService) GetTransfers(appContext context.Context) (transfers []entity.Transfer, err error) {

	log.Info("GetTransfers: Process Started")
	defer log.Info("GetTransfers: Process Finished")

	loggedAccountUUID := appContext.Value(auth.AccountUUIDKey)

	account, err := s.svc.dm.MySQL().Account().GetAccountByUUID(loggedAccountUUID.(string))
	if err != nil {
		log.Error("GetTransfers: ", err)
		return transfers, err
	}

	trasnfersMade, err := s.svc.dm.MySQL().Account().GetTransfersByAccountID(account.ID, true)
	if err != nil {
		log.Error("GetTransfers: ", err)
		return transfers, err
	}
	transfers = append(transfers, trasnfersMade...)

	trasnfersReceived, err := s.svc.dm.MySQL().Account().GetTransfersByAccountID(account.ID, false)
	if err != nil {
		log.Error("GetTransfers: ", err)
		return transfers, err
	}
	transfers = append(transfers, trasnfersReceived...)

	return transfers, nil
}
