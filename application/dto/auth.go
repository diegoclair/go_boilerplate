package dto

import (
	"time"

	"github.com/diegoclair/go_boilerplate/util/number"
	"github.com/diegoclair/go_utils/validator"
)

type Session struct {
	SessionID             int64
	SessionUUID           string `validate:"required,uuid"`
	AccountID             int64  `validate:"required"`
	RefreshToken          string `validate:"required"`
	UserAgent             string
	ClientIP              string
	IsBlocked             bool
	RefreshTokenExpiredAt time.Time
}

func (s *Session) Validate(v validator.Validator) error {
	err := v.ValidateStruct(s)
	if err != nil {
		return err
	}

	return nil
}

type LoginInput struct {
	CPF      string `validate:"required,cpf"`
	Password string `validate:"required,min=8"`
}

// Validate validate the input
func (l *LoginInput) Validate(v validator.Validator) error {
	l.CPF = number.CleanNumber(l.CPF)

	err := v.ValidateStruct(l)
	if err != nil {
		return err
	}

	return nil
}
