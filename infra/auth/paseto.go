package auth

import (
	"fmt"

	"time"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/labstack/gommon/log"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type pasetoAuth struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func newPasetoAuth(symmetricKey string) (AuthToken, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &pasetoAuth{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (a *pasetoAuth) CreateAccessToken(accountUUID, sessionUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, sessionUUID, accessTokenDurationTime)
}

func (a *pasetoAuth) CreateRefreshToken(accountUUID, sessionUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, sessionUUID, refreshTokenDurationTime)
}

func (a *pasetoAuth) VerifyToken(token string) (*tokenPayload, error) {
	payload := &tokenPayload{}

	err := a.paseto.Decrypt(token, a.symmetricKey, payload, nil)
	if err != nil {
		log.Error("VerifyToken: error to decrypt token: ", err)
		return nil, resterrors.NewUnauthorizedError(errInvalidToken.Error())
	}

	return payload, payload.Valid()
}

func (a *pasetoAuth) createToken(accountUUID, sessionUUID string, duration time.Duration) (string, *tokenPayload, error) {
	payload := newPayload(accountUUID, sessionUUID, duration)
	token, err := a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		log.Error("createToken: error to encrypt token: ", err)
		return token, payload, resterrors.NewUnauthorizedError(err.Error())
	}
	return token, payload, nil
}
