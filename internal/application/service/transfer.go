package service

import (
	"context"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/internal/domain/errcodes"
	"github.com/diegoclair/logger"
	"github.com/diegoclair/go_utils/validator"
	"github.com/google/uuid"
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
	transfer, err := input.ToEntityValidate(ctx, s.validator)
	if err != nil {
		s.log.Error(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	fromAccount, err := s.accountSvc.GetLoggedAccount(ctx)
	if err != nil {
		s.log.Error(ctx, "error to get logged account", logger.Err(err))
		return err
	}

	if !fromAccount.HasSufficientFunds(transfer.Amount) {
		return errcodes.ErrInsufficientFunds
	}

	destAccount, err := s.dm.Account().GetAccountByUUID(ctx, transfer.AccountDestinationUUID)
	if err != nil {
		if apperr.IsNotFound(err) {
			s.log.Error(ctx, "destination account not found", logger.Err(err))
			return errcodes.ErrInvalidDestinationAccount
		}
		s.log.Error(ctx, "error to get destination account by uuid", logger.Err(err))
		return err
	}

	if fromAccount.ID == destAccount.ID {
		return errcodes.ErrSelfTransfer
	}

	transfer.TransferUUID = uuid.Must(uuid.NewV7()).String()

	return s.dm.WithTransaction(ctx, func(tx contract.DataManager) error {

		_, err = tx.Account().AddTransfer(ctx, transfer.TransferUUID, fromAccount.ID, destAccount.ID, transfer.Amount)
		if err != nil {
			s.log.Error(ctx, "error to add transfer", logger.Err(err))
			return err
		}

		fromAccount.SubtractBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, fromAccount.ID, fromAccount.Balance)
		if err != nil {
			s.log.Error(ctx, "error to update origin account balance", logger.Err(err))
			return err
		}

		destAccount.AddBalance(transfer.Amount)

		err = tx.Account().UpdateAccountBalance(ctx, destAccount.ID, destAccount.Balance)
		if err != nil {
			s.log.Error(ctx, "error to update destination account balance", logger.Err(err))
			return err
		}
		return nil
	})
}

func (s *transferService) GetTransfers(ctx context.Context, take, skip int64) (transfers []entity.Transfer, totalRecords int64, err error) {
	accountID, err := s.accountSvc.GetLoggedAccountID(ctx)
	if err != nil {
		s.log.Error(ctx, "error to get logged account", logger.Err(err))
		return transfers, totalRecords, err
	}

	madeTransfers, madeTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, true)
	if err != nil {
		s.log.Error(ctx, "error to get made transfers", logger.Err(err))
		return transfers, totalRecords, err
	}

	transfers = append(transfers, madeTransfers...)

	receivedTransfers, receivedTotalRecords, err := s.dm.Account().GetTransfersByAccountID(ctx, accountID, take, skip, false)
	if err != nil {
		s.log.Error(ctx, "error to get received transfers", logger.Err(err))
		return transfers, totalRecords, err
	}

	transfers = append(transfers, receivedTransfers...)
	totalRecords = madeTotalRecords + receivedTotalRecords

	return transfers, totalRecords, nil
}
