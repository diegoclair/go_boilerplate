package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_utils/logger"
)

const (
	minSecretKeySize = 32
)

var (
	accessTokenDurationTime  time.Duration
	refreshTokenDurationTime time.Duration
)

var (
	errExpiredToken      = errors.New("token has expired")
	errInvalidToken      = errors.New("token is invalid")
	errInvalidPrivateKey = fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
)

func NewAuthToken(accessTokenDuration, refreshTokenDuration time.Duration, pasetoSymmetricKey string, log logger.Logger) (contract.AuthToken, error) {
	accessTokenDurationTime = accessTokenDuration
	refreshTokenDurationTime = refreshTokenDuration

	return newPasetoAuth(pasetoSymmetricKey, log)
}
