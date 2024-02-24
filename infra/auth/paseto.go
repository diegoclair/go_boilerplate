package auth

import (
	"context"
	"strings"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type pasetoAuth struct {
	paseto       *paseto.V2
	symmetricKey []byte
	log          logger.Logger
}

func newPasetoAuth(symmetricKey string, log logger.Logger) (AuthToken, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, errInvalidPrivateKey
	}

	return &pasetoAuth{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		log:          log,
	}, nil
}

func (a *pasetoAuth) CreateAccessToken(ctx context.Context, accountUUID, sessionUUID string) (tokenString string, payload *TokenPayload, err error) {
	return a.createToken(ctx, accountUUID, sessionUUID, accessTokenDurationTime)
}

func (a *pasetoAuth) CreateRefreshToken(ctx context.Context, accountUUID, sessionUUID string) (tokenString string, payload *TokenPayload, err error) {
	return a.createToken(ctx, accountUUID, sessionUUID, refreshTokenDurationTime)
}

func (a *pasetoAuth) VerifyToken(ctx context.Context, token string) (*TokenPayload, error) {
	if strings.TrimSpace(token) == "" {
		return nil, resterrors.NewUnauthorizedError(errInvalidToken.Error())
	}

	payload := &TokenPayload{}

	err := a.paseto.Decrypt(token, a.symmetricKey, payload, nil)
	if err != nil {
		a.log.Errorf(ctx, "error to decrypt token: %v", err)
		return nil, resterrors.NewUnauthorizedError(errInvalidToken.Error())
	}

	err = payload.Valid()
	if err != nil {
		a.log.Errorf(ctx, "error to validate token: %v", err)
		return nil, resterrors.NewUnauthorizedError(err.Error())
	}

	return payload, payload.Valid()
}

func (a *pasetoAuth) createToken(ctx context.Context, accountUUID, sessionUUID string, duration time.Duration) (tokenString string, payload *TokenPayload, err error) {
	payload = newPayload(accountUUID, sessionUUID, duration)

	tokenString, err = a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		a.log.Errorf(ctx, "error to encrypt token: %v", err)
		return tokenString, payload, resterrors.NewUnauthorizedError(err.Error())
	}

	return tokenString, payload, nil
}
