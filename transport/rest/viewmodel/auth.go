package viewmodel

import (
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
)

// validate tags are necessary to generate swagger correctly

type Login struct {
	CPF      string `json:"cpf" validate:"required,min=11,max=11"`
	Password string `json:"password" validate:"required,min=8"`
}

func (l *Login) ToDto() dto.LoginInput {
	return dto.LoginInput{
		CPF:      l.CPF,
		Password: l.Password,
	}
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

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
