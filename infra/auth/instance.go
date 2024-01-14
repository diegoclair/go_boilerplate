package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_boilerplate/infra/logger"
)

type AuthToken interface {
	CreateAccessToken(ctx context.Context, accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error)
	CreateRefreshToken(ctx context.Context, accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error)
	VerifyToken(ctx context.Context, token string) (*tokenPayload, error)
}

const (
	tokenTypeJWT     = "jwt"
	tokenTypePaseto  = "paseto"
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

func NewAuthToken(cfg config.AuthConfig, log logger.Logger) (AuthToken, error) {
	accessTokenDurationTime = cfg.AccessTokenDuration
	refreshTokenDurationTime = cfg.RefreshTokenDuration

	if cfg.AccessTokenType == tokenTypeJWT {
		return newJwtAuth(cfg.JWTPrivateKey, log)
	}
	return newPasetoAuth(cfg.PasetoSymmetricKey, log)
}
