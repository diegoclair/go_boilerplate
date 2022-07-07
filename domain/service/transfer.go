package service

import (
	"context"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
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

func (s *transferService) CreateTransfer(ctx context.Context, transfer entity.Transfer) (err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	loggedAccountUUID := ctx.Value(auth.AccountUUIDKey)
	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, loggedAccountUUID.(string))
	if err != nil {
		log.Error(err)
		return err
	}

	if account.Balance < transfer.Amount {
		return resterrors.NewConflictError("Your account don't have sufficient funds to do this operation")
	}

	destAccount, err := s.svc.dm.Account().GetAccountByUUID(ctx, transfer.AccountDestinationUUID)
	if err != nil {
		log.Error(err)
		return err
	}

	transfer.AccountDestinationID = destAccount.ID
	transfer.AccountOriginID = account.ID
	transfer.TransferUUID = uuid.NewV4().String()

	tx, err := s.svc.dm.Begin()
	if err != nil {
		log.Error("error to get db transaction", err)
		return err
	}
	defer tx.Rollback()

	err = tx.Account().AddTransfer(ctx, transfer)
	if err != nil {
		log.Error(err)
		return err
	}

	originBalance := account.Balance - transfer.Amount
	err = tx.Account().UpdateAccountBalance(ctx, account.ID, originBalance)
	if err != nil {
		log.Error(err)
		return err
	}

	destBalance := destAccount.Balance + transfer.Amount
	err = tx.Account().UpdateAccountBalance(ctx, destAccount.ID, destBalance)
	if err != nil {
		log.Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *transferService) GetTransfers(ctx context.Context) (transfers []entity.Transfer, err error) {

	ctx, log := s.svc.log.NewSessionLogger(ctx)
	log.Info("Process Started")
	defer log.Info("Process Finished")

	loggedAccountUUID := ctx.Value(auth.AccountUUIDKey)

	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, loggedAccountUUID.(string))
	if err != nil {
		log.Error(err)
		return transfers, err
	}

	trasnfersMade, err := s.svc.dm.Account().GetTransfersByAccountID(ctx, account.ID, true)
	if err != nil {
		log.Error(err)
		return transfers, err
	}
	transfers = append(transfers, trasnfersMade...)

	trasnfersReceived, err := s.svc.dm.Account().GetTransfersByAccountID(ctx, account.ID, false)
	if err != nil {
		log.Error(err)
		return transfers, err
	}
	transfers = append(transfers, trasnfersReceived...)

	return transfers, nil
}
