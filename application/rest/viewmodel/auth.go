package viewmodel

import (
	"time"

	"github.com/diegoclair/go_utils-lib/v2/validator"
)

type Login struct {
	CPF      string `json:"cpf,omitempty" validate:"required,min=11,max=11,cpf"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

func (l *Login) Validate(structValidator validator.Validator) error {
	err := structValidator.ValidateStruct(l)
	if err != nil {
		return err
	}

	return nil
}

type LoginResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (t *RefreshTokenRequest) Validate(structValidator validator.Validator) error {
	err := structValidator.ValidateStruct(t)
	if err != nil {
		return err
	}

	return nil
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
