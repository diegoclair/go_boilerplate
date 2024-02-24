package dto

import (
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/util/number"
	"github.com/diegoclair/go_utils/validator"
)

type AccountInput struct {
	Name     string `validate:"required,min=3"`
	CPF      string `validate:"required,cpf"`
	Password string `validate:"required,min=8"`
}

// ToEntityValidate validate the input and return the entity
func (a *AccountInput) ToEntityValidate(validator validator.Validator) (account entity.Account, err error) {
	a.CPF = number.CleanNumber(a.CPF)

	err = validator.ValidateStruct(a)
	if err != nil {
		return account, err
	}

	account = entity.Account{
		Name:     a.Name,
		CPF:      a.CPF,
		Password: a.Password,
	}

	return account, nil
}

type AddBalanceInput struct {
	AccountUUID string  `validate:"required,uuid"`
	Amount      float64 `validate:"required,gt=0"`
}

// Validate validate the input
func (a *AddBalanceInput) Validate(validator validator.Validator) error {
	err := validator.ValidateStruct(a)
	if err != nil {
		return err
	}

	return nil
}
