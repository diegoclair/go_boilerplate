package auth

import (
	"fmt"

	"time"

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

func (a *pasetoAuth) CreateToken(accountUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, accessTokenDurationTime)
}

func (a *pasetoAuth) CreateRefreshToken(accountUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, refreshTokenDurationTime)
}

func (a *pasetoAuth) VerifyToken(token string) (*tokenPayload, error) {
	payload := &tokenPayload{}

	err := a.paseto.Decrypt(token, a.symmetricKey, payload, nil)
	if err != nil {
		log.Error("error to decrypt token: ", err)
		return nil, errInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		log.Error("token not valid: ", err)
		return nil, errInvalidToken
	}

	return payload, nil
}

func (a *pasetoAuth) createToken(accountUUID string, duration time.Duration) (string, *tokenPayload, error) {
	payload := newPayload(accountUUID, duration)
	token, err := a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		log.Error("error to encrypt token: ", err)
	}
	return token, payload, err
}
