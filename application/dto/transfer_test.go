package dto

import (
	"context"
	"testing"

	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferInput_ToEntityValidate(t *testing.T) {
	ctx := context.Background()
	v, err := validator.NewValidator()
	require.NoError(t, err)

	type fields struct {
		AccountDestinationUUID string
		Amount                 float64
	}

	tests := []struct {
		name         string
		fields       fields
		wantTransfer entity.Transfer
		wantErr      bool
	}{
		{
			name: "Should return transfer entity without error",
			fields: fields{
				AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				Amount:                 5.0,
			},

			wantTransfer: entity.Transfer{
				AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				Amount:                 5.0,
			},
			wantErr: false,
		},
		{
			name: "Should return error if account destination uuid is empty",
			fields: fields{
				AccountDestinationUUID: "",
				Amount:                 5.0,
			},
			wantTransfer: entity.Transfer{},
			wantErr:      true,
		},
		{
			name: "Should return error if amount is empty",
			fields: fields{
				AccountDestinationUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				Amount:                 0,
			},
			wantTransfer: entity.Transfer{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tInput := &TransferInput{
				AccountDestinationUUID: tt.fields.AccountDestinationUUID,
				Amount:                 tt.fields.Amount,
			}

			gotTransfer, err := tInput.ToEntityValidate(ctx, v)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransferInput.ToEntityValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantTransfer, gotTransfer)
		})
	}
}
