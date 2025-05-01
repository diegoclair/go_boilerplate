package service

import (
	"context"

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
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	fromAccount, err := s.accountSvc.GetLoggedAccount(ctx)
	if err != nil {
		s.log.Errorw(ctx, "error to get logged account", logger.Err(err))
		return err
	}

	if !fromAccount.HasSufficientFunds(transfer.Amount) {
		return resterrors.NewConflictError("Your account don't have sufficient funds to do this operation")
	}

	destAccount, err := s.dm.Account().GetAccountByUUID(ctx, transfer.AccountDestinationUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			s.log.Errorw(ctx, "error to get destination account by uuid", logger.Err(err))
			return resterrors.NewNotFoundError("Invalid destination account")
		}
		s.log.Errorw(ctx, "error to get destination account by uuid", logger.Err(err))
		return err
	}

	if fromAccount.ID == destAccount.ID {
		return resterrors.NewConflictError("You can't transfer to yourself")
	}

	transfer.TransferUUID = uuid.NewV4().String()

	return s.dm.WithTransaction(ctx, func(tx contract.DataManager) error {

		_, err = tx.Account().AddTransfer(ctx, transfer.TransferUUID, fromAccount.ID, destAccount.ID, transfer.Amount)
		if err != nil {
			s.log.Errorw(ctx, "error to add transfer", logger.Err(err))
			return err
		}

		fromAccount.SubtractBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, fromAccount.ID, fromAccount.Balance)
		if err != nil {
			s.log.Errorw(ctx, "error to update origin account balance", logger.Err(err))
			return err
		}

		destAccount.AddBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, destAccount.ID, destAccount.Balance)
		if err != nil {
			s.log.Errorw(ctx, "error to update destination account balance", logger.Err(err))
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
		s.log.Errorw(ctx, "error to get logged account", logger.Err(err))
		return transfers, totalRecords, err
	}

	madeTransfers, madeTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, true)
	if err != nil {
		s.log.Errorw(ctx, "error to get made transfers", logger.Err(err))
		return transfers, totalRecords, err
	}

	transfers = append(transfers, madeTransfers...)

	receivedTransfers, receivedTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, false)
	if err != nil {
		s.log.Errorw(ctx, "error to get received transfers", logger.Err(err))
		return transfers, totalRecords, err
	}

	transfers = append(transfers, receivedTransfers...)
	totalRecords = madeTotalRecords + receivedTotalRecords

	return transfers, totalRecords, nil
}
