package viewmodel

import (
	"time"

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

type LoginResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (t *RefreshTokenRequest) Validate() error {

	err := validstruct.ValidateStruct(t)
	if err != nil {
		return err
	}

	return nil
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
