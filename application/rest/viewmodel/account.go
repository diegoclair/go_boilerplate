package viewmodel

import (
	"time"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils-lib/v2/validator"
)

type AddAccount struct {
	Name     string `json:"name,omitempty" validate:"required,min=3"`
	CPF      string `json:"cpf,omitempty" validate:"required,min=11,max=11,cpf"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

func (a *AddAccount) Validate(strucValidator validator.Validator) error {
	err := strucValidator.ValidateStruct(a)
	if err != nil {
		return err
	}

	return nil
}

type Account struct {
	UUID      string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CPF       string    `json:"cpf,omitempty"`
	Balance   float64   `json:"balance"`
	CreatedAT time.Time `json:"create_at,omitempty"`
}

func (a *Account) FillFromEntity(account entity.Account) {
	a.UUID = account.UUID
	a.Name = account.Name
	a.CPF = account.CPF
	a.Balance = account.Balance
	a.CreatedAT = account.CreatedAT
}

type AddBalance struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

func (a *AddBalance) Validate(strucValidator validator.Validator) error {
	err := strucValidator.ValidateStruct(a)
	if err != nil {
		return err
	}

	return nil
}
