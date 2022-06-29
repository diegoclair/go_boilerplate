package auth

import (
	"time"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/gommon/log"
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
		log.Error(errExpiredToken)
		return resterrors.NewUnauthorizedError(errExpiredToken.Error())
	}
	return nil
}
