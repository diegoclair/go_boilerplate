package viewmodel

import (
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/entity"
)

// validate tags are necessary to generate swagger correctly

type TransferReq struct {
	AccountDestinationUUID string  `json:"account_destination_id" validate:"required,uuid"`
	Amount                 float64 `json:"amount" validate:"required,gt=0"`
}

func (t *TransferReq) ToDto() dto.TransferInput {
	return dto.TransferInput{
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

func (t *TransferResp) FillFromEntity(transfer entity.Transfer) {
	t.TransferUUID = transfer.TransferUUID
	t.AccountOriginUUID = transfer.AccountOriginUUID
	t.AccountDestinationUUID = transfer.AccountDestinationUUID
	t.Amount = transfer.Amount
	t.CreateAt = transfer.CreatedAt
}
