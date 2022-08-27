package auth

import (
	"errors"
	"time"

	"github.com/diegoclair/go_boilerplate/util/config"
)

// TODO: Create unit tests for this package
type AuthToken interface {
	CreateAccessToken(accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error)
	CreateRefreshToken(accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error)
	VerifyToken(token string) (*tokenPayload, error)
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
	errExpiredToken = errors.New("token has expired")
	errInvalidToken = errors.New("token is invalid")
)

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	AccountUUIDKey Key = "AccountUUID"
	TokenKey       Key = "user-token"
	SessionKey     Key = "Session"
)

func NewAuthToken(cfg config.AuthConfig) (AuthToken, error) {
	accessTokenDurationTime = cfg.AccessTokenDuration
	refreshTokenDurationTime = cfg.RefreshTokenDuration

	if cfg.AccessTokenType == tokenTypePaseto {
		return newPasetoAuth(cfg.PasetoSymmetricKey)
	}
	return newJwtAuth(cfg.JWTPrivateKey)

}
