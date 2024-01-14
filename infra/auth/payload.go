package auth

import (
	"time"
)

type tokenPayload struct {
	AccountUUID  string
	SessionUUID  string
	RefreshToken string
	IssuedAt     time.Time
	ExpiredAt    time.Time
}

func newPayload(accountUUID, sessionUUID string, duration time.Duration) *tokenPayload {
	return &tokenPayload{
		SessionUUID: sessionUUID,
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
