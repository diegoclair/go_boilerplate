package dto

import (
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_utils/validator"
	"golang.org/x/net/context"
)

type TransferInput struct {
	AccountDestinationUUID string  `validate:"required,uuid"`
	Amount                 float64 `validate:"required"`
}

// ToEntityValidate validate the input and return the entity
func (t *TransferInput) ToEntityValidate(ctx context.Context, v validator.Validator) (transfer entity.Transfer, err error) {
	err = v.ValidateStruct(ctx, t)
	if err != nil {
		return transfer, err
	}

	return entity.Transfer{
		AccountDestinationUUID: t.AccountDestinationUUID,
		Amount:                 t.Amount,
	}, nil
}
