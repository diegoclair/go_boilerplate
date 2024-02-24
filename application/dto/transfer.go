package dto

import (
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils/validator"
)

type TransferInput struct {
	AccountDestinationUUID string  `validate:"required,uuid"`
	Amount                 float64 `validate:"required"`
}

// ToEntityValidate validate the input and return the entity
func (t *TransferInput) ToEntityValidate(v validator.Validator) (transfer entity.Transfer, err error) {
	err = v.ValidateStruct(t)
	if err != nil {
		return transfer, err
	}

	return entity.Transfer{
		AccountDestinationUUID: t.AccountDestinationUUID,
		Amount:                 t.Amount,
	}, nil
}
