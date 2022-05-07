package auth

import (
	"time"

	"github.com/twinj/uuid"
)

type tokenPayload struct {
	ID          uuid.UUID `json:"id"`
	AccountUUID string    `json:"account_uuid"`
	SessionID   string    `json:"session_id"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func newPayload(accountUUID string, duration time.Duration) *tokenPayload {
	return &tokenPayload{
		ID:          uuid.NewV4(),
		AccountUUID: accountUUID,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
}

// Valid checks if the token payload is valid or not
func (p *tokenPayload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errExpiredToken
	}
	return nil
}
