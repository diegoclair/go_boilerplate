package viewmodel

import (
	"time"

	"github.com/diegoclair/go_utils-lib/v2/validstruct"
)

type Transfer struct {
	TransferUUID           string    `json:"id"`
	AccountOriginUUID      string    `json:"account_origin_id,omitempty"`
	AccountDestinationUUID string    `json:"account_destination_id,omitempty" validate:"required"`
	Amount                 float64   `json:"amount,omitempty" validate:"required"`
	CreateAt               time.Time `json:"create_at,omitempty"`
}

func (t *Transfer) Validate() error {
	return validstruct.ValidateStruct(t)
}
