package viewmodel

import (
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
)

type AddAccount struct {
	Name     string `json:"name,omitempty"`
	CPF      string `json:"cpf,omitempty"`
	Password string `json:"password,omitempty"`
}

func (a *AddAccount) ToDto() dto.AccountInput {
	return dto.AccountInput{
		Name:     a.Name,
		CPF:      a.CPF,
		Password: a.Password,
	}
}

type AccountResponse struct {
	UUID      string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CPF       string    `json:"cpf,omitempty"`
	Balance   float64   `json:"balance"`
	CreatedAT time.Time `json:"create_at,omitempty"`
}

func (a *AccountResponse) FillFromEntity(account entity.Account) {
	a.UUID = account.UUID
	a.Name = account.Name
	a.CPF = account.CPF
	a.Balance = account.Balance
	a.CreatedAT = account.CreatedAT
}

type AddBalance struct {
	Amount float64 `json:"amount"`
}

func (a *AddBalance) ToDto(accountUUID string) dto.AddBalanceInput {
	return dto.AddBalanceInput{
		AccountUUID: accountUUID,
		Amount:      a.Amount,
	}
}
