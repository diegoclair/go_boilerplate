package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/logger"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type pasetoAuth struct {
	paseto       *paseto.V2
	symmetricKey []byte
	log          logger.Logger
}

func newPasetoAuth(symmetricKey string, log logger.Logger) (*pasetoAuth, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, errInvalidPrivateKey
	}

	return &pasetoAuth{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		log:          log,
	}, nil
}

func (p *pasetoAuth) CreateAccessToken(ctx context.Context, input contract.TokenPayloadInput) (tokenString string, resp contract.TokenPayload, err error) {
	tokenString, payload, err := p.createToken(ctx, fromContractTokenPayloadInput(input), accessTokenDurationTime)
	if err != nil {
		return tokenString, resp, err
	}

	return tokenString, payload.toContract(), nil
}

func (p *pasetoAuth) CreateRefreshToken(ctx context.Context, input contract.TokenPayloadInput) (tokenString string, resp contract.TokenPayload, err error) {
	tokenString, payload, err := p.createToken(ctx, fromContractTokenPayloadInput(input), refreshTokenDurationTime)
	if err != nil {
		return tokenString, resp, err
	}

	return tokenString, payload.toContract(), nil
}

func (p *pasetoAuth) VerifyToken(ctx context.Context, token string) (resp contract.TokenPayload, err error) {
	if strings.TrimSpace(token) == "" {
		return resp, apperr.ErrTokenInvalid
	}

	payload := &tokenPayload{}

	err = p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		p.log.Error(ctx, "error to decrypt token", logger.Err(err))
		return resp, apperr.ErrTokenInvalid
	}

	err = payload.Valid()
	if err != nil {
		if errors.Is(err, errExpiredToken) {
			p.log.Warn(ctx, "token has expired")
			return resp, apperr.ErrTokenExpired
		}
		p.log.Error(ctx, "error to validate token", logger.Err(err))
		return resp, apperr.ErrTokenInvalid
	}

	return payload.toContract(), nil
}

func (a *pasetoAuth) createToken(ctx context.Context, input tokenPayloadInput, duration time.Duration) (tokenString string, payload *tokenPayload, err error) {
	payload = newPayload(input, duration)

	tokenString, err = a.paseto.Encrypt(a.symmetricKey, payload, nil)
	if err != nil {
		a.log.Error(ctx, "error to encrypt token", logger.Err(err))
		return tokenString, payload, apperr.ErrInternal.Wrap(err)
	}

	return tokenString, payload, nil
}
