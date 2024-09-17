package service

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/twinj/uuid"
)

type transferService struct {
	accountSvc contract.AccountApp
	dm         contract.DataManager
	log        logger.Logger
	validator  validator.Validator
}

func newTransferService(infra domain.Infrastructure, accountSvc contract.AccountApp) *transferService {
	return &transferService{
		accountSvc: accountSvc,
		dm:         infra.DataManager(),
		log:        infra.Logger(),
		validator:  infra.Validator(),
	}
}

func (s *transferService) CreateTransfer(ctx context.Context, input dto.TransferInput) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	transfer, err := input.ToEntityValidate(ctx, s.validator)
	if err != nil {
		s.log.Errorf(ctx, "error or invalid input: %s", err.Error())
		return err
	}

	fromAccount, err := s.accountSvc.GetLoggedAccount(ctx)
	if err != nil {
		s.log.Errorf(ctx, "error to get logged account: %s", err.Error())
		return err
	}

	if !fromAccount.HasSufficientFunds(transfer.Amount) {
		return resterrors.NewConflictError("Your account don't have sufficient funds to do this operation")
	}

	destAccount, err := s.dm.Account().GetAccountByUUID(ctx, transfer.AccountDestinationUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			s.log.Errorf(ctx, "error to get destination account by uuid: %s", err.Error())
			return resterrors.NewNotFoundError("Invalid destination account")
		}
		s.log.Errorf(ctx, "error to get destination account by uuid: %s", err.Error())
		return err
	}

	if fromAccount.ID == destAccount.ID {
		return resterrors.NewConflictError("You can't transfer to yourself")
	}

	transfer.TransferUUID = uuid.NewV4().String()

	return s.dm.WithTransaction(ctx, func(tx contract.DataManager) error {

		_, err = tx.Account().AddTransfer(ctx, transfer.TransferUUID, fromAccount.ID, destAccount.ID, transfer.Amount)
		if err != nil {
			s.log.Errorf(ctx, "error to add transfer: %s", err.Error())
			return err
		}

		fromAccount.SubtractBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, fromAccount.ID, fromAccount.Balance)
		if err != nil {
			s.log.Errorf(ctx, "error to update origin account balance: %s", err.Error())
			return err
		}

		destAccount.AddBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, destAccount.ID, destAccount.Balance)
		if err != nil {
			s.log.Errorf(ctx, "error to update destination account balance: %s", err.Error())
			return err
		}
		return nil
	})
}

func (s *transferService) GetTransfers(ctx context.Context, take, skip int64) (transfers []entity.Transfer, totalRecords int64, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	accountID, err := s.accountSvc.GetLoggedAccountID(ctx)
	if err != nil {
		s.log.Errorf(ctx, "error to get logged account: %s", err.Error())
		return transfers, totalRecords, err
	}

	madeTransfers, madeTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, true)
	if err != nil {
		s.log.Errorf(ctx, "error to get made transfers: %s", err.Error())
		return transfers, totalRecords, err
	}

	transfers = append(transfers, madeTransfers...)

	receivedTransfers, receivedTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, false)
	if err != nil {
		s.log.Errorf(ctx, "error to get received transfers: %s", err.Error())
		return transfers, totalRecords, err
	}

	transfers = append(transfers, receivedTransfers...)
	totalRecords = madeTotalRecords + receivedTotalRecords

	return transfers, totalRecords, nil
}
