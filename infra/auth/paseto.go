package auth

import (
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
		return nil, errInvalidPrivateKey
	}

	return &pasetoAuth{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (a *pasetoAuth) CreateAccessToken(accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error) {
	return a.createToken(accountUUID, sessionUUID, accessTokenDurationTime)
}

func (a *pasetoAuth) CreateRefreshToken(accountUUID, sessionUUID string) (tokenString string, payload *tokenPayload, err error) {
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

func (a *pasetoAuth) createToken(accountUUID, sessionUUID string, duration time.Duration) (tokenString string, payload *tokenPayload, err error) {
	payload = newPayload(accountUUID, sessionUUID, duration)
	tokenString, err = a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		log.Error("createToken: error to encrypt token: ", err)
		return tokenString, payload, resterrors.NewUnauthorizedError(err.Error())
	}
	return tokenString, payload, nil
}
