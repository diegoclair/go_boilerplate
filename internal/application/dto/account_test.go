package dto

import (
	"context"
	"testing"

	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_utils/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountInput_ToEntityValidate(t *testing.T) {
	ctx := context.Background()
	v, err := validator.NewValidator()
	require.NoError(t, err)

	type fields struct {
		Name     string
		CPF      string
		Password string
	}

	tests := []struct {
		name       string
		fields     fields
		wantEntity entity.Account
		wantErr    bool
	}{
		{
			name: "Should return account entity without error",
			fields: fields{
				Name:     "John Doe",
				CPF:      "01234567890",
				Password: "12345678",
			},
			wantEntity: entity.Account{
				Name:     "John Doe",
				CPF:      "01234567890",
				Password: "12345678",
			},
			wantErr: false,
		},
		{
			name: "Should return error if name is empty",
			fields: fields{
				Name:     "",
				CPF:      "01234567890",
				Password: "12345678",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
		{
			name: "Should return error if name is less than 3 characters",
			fields: fields{
				Name:     "Jo",
				CPF:      "01234567890",
				Password: "12345678",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
		{
			name: "Should return error if cpf is empty",
			fields: fields{
				Name:     "John Doe",
				CPF:      "",
				Password: "12345678",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
		{
			name: "Should return error if cpf is invalid",
			fields: fields{
				Name:     "John Doe",
				CPF:      "1234567890",
				Password: "12345678",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
		{
			name: "Should return error if password is empty",
			fields: fields{
				Name:     "John Doe",
				CPF:      "01234567890",
				Password: "",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
		{
			name: "Should return error if password is less than 8 characters",
			fields: fields{
				Name:     "John Doe",
				CPF:      "01234567890",
				Password: "1234567",
			},
			wantEntity: entity.Account{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aInput := &AccountInput{
				Name:     tt.fields.Name,
				CPF:      tt.fields.CPF,
				Password: tt.fields.Password,
			}

			gotEntity, err := aInput.ToEntityValidate(ctx, v)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountInput.ToEntityValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantEntity, gotEntity)
		})
	}
}

func TestAddBalanceInput_Validate(t *testing.T) {
	ctx := context.Background()
	v, err := validator.NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  AddBalanceInput
		wantErr bool
	}{
		{
			name: "Should return without error",
			fields: AddBalanceInput{
				AccountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				Amount:      5.0,
			},
			wantErr: false,
		},
		{
			name: "Should return error if account uuid is empty",
			fields: AddBalanceInput{
				AccountUUID: "",
				Amount:      5.0,
			},
			wantErr: true,
		},
		{
			name: "Should return error if account uuid is invalid",
			fields: AddBalanceInput{
				AccountUUID: "d152a340-9a87-4d32-85ad-19df4c9934c",
				Amount:      5.0,
			},
			wantErr: true,
		},
		{
			name: "Should return error if amount is empty",
			fields: AddBalanceInput{
				AccountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				Amount:      0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.fields.Validate(ctx, v)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBalanceInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
