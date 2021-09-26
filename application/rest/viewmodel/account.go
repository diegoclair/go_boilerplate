package viewmodel

import (
	"time"

	"github.com/diegoclair/go_utils-lib/v2/validstruct"
)

type AddAccount struct {
	Name string `json:"name,omitempty" validate:"required,min=3"`
	Login
}

func (u *AddAccount) Validate() error {
	err := u.Login.Validate()
	if err != nil {
		return err
	}

	err = validstruct.ValidateStruct(u)
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
