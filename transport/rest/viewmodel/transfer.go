package viewmodel

import (
	"time"

	"github.com/diegoclair/go_boilerplate/domain/transfer"
	"github.com/diegoclair/go_utils/validator"
)

type TransferReq struct {
	AccountDestinationUUID string  `json:"account_destination_id,omitempty" validate:"required"`
	Amount                 float64 `json:"amount,omitempty" validate:"required"`
}

func (t *TransferReq) Validate(v validator.Validator) error {
	return v.ValidateStruct(t)
}

func (t *TransferReq) ToEntity() transfer.Transfer {
	return transfer.Transfer{
		AccountDestinationUUID: t.AccountDestinationUUID,
		Amount:                 t.Amount,
	}
}

type TransferResp struct {
	TransferUUID           string    `json:"id"`
	AccountOriginUUID      string    `json:"account_origin_id,omitempty"`
	AccountDestinationUUID string    `json:"account_destination_id,omitempty"`
	Amount                 float64   `json:"amount,omitempty"`
	CreateAt               time.Time `json:"create_at,omitempty"`
}

func (t *TransferResp) FillFromEntity(transfer transfer.Transfer) {
	t.TransferUUID = transfer.TransferUUID
	t.AccountOriginUUID = transfer.AccountOriginUUID
	t.AccountDestinationUUID = transfer.AccountDestinationUUID
	t.Amount = transfer.Amount
	t.CreateAt = transfer.CreatedAt
}
