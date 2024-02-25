package auth

import (
	"time"
)

// TokenPayload represents the payload of a JWT token
type TokenPayload struct {
	AccountUUID  string
	SessionUUID  string
	RefreshToken string
	IssuedAt     time.Time
	ExpiredAt    time.Time
}

func newPayload(input TokenPayloadInput, duration time.Duration) *TokenPayload {
	return &TokenPayload{
		SessionUUID: input.SessionUUID,
		AccountUUID: input.AccountUUID,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
}

// Valid checks if the token payload is valid or not
func (p *TokenPayload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errExpiredToken
	}
	return nil
}
