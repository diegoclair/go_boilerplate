package service

import (
	"context"
	"errors"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/util/number"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/twinj/uuid"
)

type transferService struct {
	svc *service
}

func newTransferService(svc *service) contract.TransferService {
	return &transferService{
		svc: svc,
	}
}

func (s *transferService) CreateTransfer(ctx context.Context, transfer entity.Transfer) (err error) {

	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	loggedAccountUUID, ok := ctx.Value(auth.AccountUUIDKey).(string)
	if !ok {
		errMsg := "accountUUID should not be empty"
		s.svc.log.Error(ctx, errMsg)
		return errors.New(errMsg)
	}
	fromAccount, err := s.svc.dm.Account().GetAccountByUUID(ctx, loggedAccountUUID)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return err
	}

	if fromAccount.Balance < transfer.Amount {
		return resterrors.NewConflictError("Your account don't have sufficient funds to do this operation")
	}

	destAccount, err := s.svc.dm.Account().GetAccountByUUID(ctx, transfer.AccountDestinationUUID)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return err
	}

	transfer.TransferUUID = uuid.NewV4().String()

	return s.svc.dm.WithTransaction(ctx, func(tx contract.DataManager) error {

		err = tx.Account().AddTransfer(ctx, transfer.TransferUUID, fromAccount.ID, destAccount.ID, transfer.Amount)
		if err != nil {
			s.svc.log.Error(ctx, err.Error())
			return err
		}

		originBalance := number.RoundFloat(fromAccount.Balance-transfer.Amount, 2)
		err = tx.Account().UpdateAccountBalance(ctx, fromAccount.ID, originBalance)
		if err != nil {
			s.svc.log.Error(ctx, err.Error())
			return err
		}

		destBalance := number.RoundFloat(destAccount.Balance+transfer.Amount, 2)
		err = tx.Account().UpdateAccountBalance(ctx, destAccount.ID, destBalance)
		if err != nil {
			s.svc.log.Error(ctx, err.Error())
			return err
		}
		return nil
	})
}

func (s *transferService) GetTransfers(ctx context.Context) (transfers []entity.Transfer, err error) {

	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	loggedAccountUUID, ok := ctx.Value(auth.AccountUUIDKey).(string)
	if !ok {
		errMsg := "accountUUID should not be empty"
		s.svc.log.Error(ctx, errMsg)
		return transfers, errors.New(errMsg)
	}

	account, err := s.svc.dm.Account().GetAccountByUUID(ctx, loggedAccountUUID)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return transfers, err
	}

	transfersMade, err := s.svc.dm.Account().GetTransfersByAccountID(ctx, account.ID, true)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return transfers, err
	}
	transfers = append(transfers, transfersMade...)

	transfersReceived, err := s.svc.dm.Account().GetTransfersByAccountID(ctx, account.ID, false)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return transfers, err
	}
	transfers = append(transfers, transfersReceived...)

	return transfers, nil
}
