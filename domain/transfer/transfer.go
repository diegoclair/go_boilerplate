package transfer

import "time"

type Transfer struct {
	ID                     int64
	TransferUUID           string
	AccountOriginUUID      string
	AccountDestinationUUID string
	Amount                 float64
	CreatedAt              time.Time
}
