package dto

import (
	"testing"

	"github.com/diegoclair/go_utils/validator"
	"github.com/stretchr/testify/require"
)

func TestSession_ToEntityValidate(t *testing.T) {
	v, err := validator.NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  Session
		wantErr bool
	}{
		{
			name: "Valid session",
			fields: Session{
				SessionUUID:  "d152a340-9a87-4d32-85ad-19df4c9934cd",
				AccountID:    1,
				RefreshToken: "token",
			},
			wantErr: false,
		},
		{
			name: "Should return error if session uuid is empty",
			fields: Session{
				AccountID:    1,
				RefreshToken: "token",
			},
			wantErr: true,
		},
		{
			name: "Should return error if account id is empty",
			fields: Session{
				SessionUUID:  "d152a340-9a87-4d32-85ad-19df4c9934cd",
				RefreshToken: "token",
			},
			wantErr: true,
		},
		{
			name: "Should return error if refresh token is empty",
			fields: Session{
				SessionUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				AccountID:   1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.fields.Validate(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLoginInput_Validate(t *testing.T) {
	v, err := validator.NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  LoginInput
		wantErr bool
	}{
		{
			name: "Valid login",
			fields: LoginInput{
				CPF:      "01234567890",
				Password: "12345678",
			},
			wantErr: false,
		},
		{
			name: "Should return error if cpf is empty",
			fields: LoginInput{
				Password: "12345678",
			},
			wantErr: true,
		},
		{
			name: "Should return error if password is empty",
			fields: LoginInput{
				CPF: "12345678901",
			},
			wantErr: true,
		},
		{
			name: "Should return error if cpf is invalid",
			fields: LoginInput{
				CPF:      "1234567890",
				Password: "12345678",
			},
			wantErr: true,
		},
		{
			name: "Should return error if password is less than 8 characters",
			fields: LoginInput{
				CPF:      "12345678901",
				Password: "1234567",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.fields.Validate(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
