package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/chacha20poly1305"
)

type jwtAuth struct {
	jwtPrivateKey string
}

func newJwtAuth(jwtPrivateKey string) (AuthToken, error) {
	if len(jwtPrivateKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &jwtAuth{
		jwtPrivateKey: jwtPrivateKey,
	}, nil
}

func (a *jwtAuth) CreateToken(accountUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, accessTokenDurationTime)
}

func (a *jwtAuth) CreateRefreshToken(accountUUID string) (string, *tokenPayload, error) {
	return a.createToken(accountUUID, refreshTokenDurationTime)
}

func (a *jwtAuth) VerifyToken(token string) (*tokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidToken
		}
		return []byte(a.jwtPrivateKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &tokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errExpiredToken) {
			log.Error("expired token: ", err)
			return nil, errExpiredToken
		}
		return nil, errInvalidToken
	}

	payload, ok := jwtToken.Claims.(*tokenPayload)
	if !ok {
		log.Error("could not parse jwt token: ", err)
		return nil, errInvalidToken
	}
	return payload, nil
}

func (a *jwtAuth) createToken(accountUUID string, duration time.Duration) (string, *tokenPayload, error) {
	key := []byte(a.jwtPrivateKey)
	payload := newPayload(accountUUID, duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Error("error to encrypt token: ", err)
	}

	return tokenString, payload, err
}
