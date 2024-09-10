package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diegoclair/go_utils/logger"
)

type TokenPayloadInput struct {
	AccountUUID string
	SessionUUID string
}

type AuthToken interface {
	CreateAccessToken(ctx context.Context, input TokenPayloadInput) (tokenString string, payload *TokenPayload, err error)
	CreateRefreshToken(ctx context.Context, input TokenPayloadInput) (tokenString string, payload *TokenPayload, err error)
	VerifyToken(ctx context.Context, token string) (*TokenPayload, error)
}

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

func NewAuthToken(accessTokenDuration, refreshTokenDuration time.Duration, pasetoSymmetricKey string, log logger.Logger) (AuthToken, error) {
	accessTokenDurationTime = accessTokenDuration
	refreshTokenDurationTime = refreshTokenDuration

	return newPasetoAuth(pasetoSymmetricKey, log)
}
