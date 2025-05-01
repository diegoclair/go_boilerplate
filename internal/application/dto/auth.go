package dto

import (
	"time"

	"github.com/diegoclair/go_boilerplate/util/number"
	"github.com/diegoclair/go_utils/validator"
	"golang.org/x/net/context"
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

func (s *Session) Validate(ctx context.Context, v validator.Validator) error {
	return v.ValidateStruct(ctx, s)
}

type LoginInput struct {
	CPF      string `validate:"required,cpf"`
	Password string `validate:"required,min=8"`
}

// Validate validate the input
func (l *LoginInput) Validate(ctx context.Context, v validator.Validator) error {
	l.CPF = number.CleanNumber(l.CPF)
	return v.ValidateStruct(ctx, l)
}
