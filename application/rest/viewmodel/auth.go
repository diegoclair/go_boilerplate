package viewmodel

import (
	"github.com/diegoclair/go-boilerplate/util/validator"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/diegoclair/go_utils-lib/v2/validstruct"
)

type Login struct {
	CPF    string `json:"cpf,omitempty" validate:"required,min=11,max=11"`
	Secret string `json:"secret,omitempty" validate:"required,min=8"`
}

func (l *Login) Validate() error {

	l.CPF = validator.CleanNumber(l.CPF)
	err := validstruct.ValidateStruct(l)
	if err != nil {
		return err
	}

	validDocument := validator.IsValidCPF(l.CPF)
	if !validDocument {
		return resterrors.NewUnprocessableEntity("Invalid cpf document")
	}

	return nil
}

type AuthResponse struct {
	Token      string `json:"token"`
	ValidTime  int64  `json:"valid_time"`
	ServerTime int64  `json:"server_time"`
}
