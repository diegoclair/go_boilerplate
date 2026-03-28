package auth

import (
	"context"
	"errors"
	"strings"

	"aidanwoods.dev/go-paseto"
	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/logger"
)

const claimData = "data"

type pasetoAuth struct {
	key    paseto.V4SymmetricKey
	parser paseto.Parser
	log    logger.Logger
}

func newPasetoAuth(symmetricKey string, log logger.Logger) (*pasetoAuth, error) {
	if len(symmetricKey) != minSecretKeySize {
		return nil, errInvalidPrivateKey
	}

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, errInvalidPrivateKey
	}

	return &pasetoAuth{
		key:    key,
		parser: paseto.NewParser(),
		log:    log,
	}, nil
}

func (p *pasetoAuth) CreateAccessToken(ctx context.Context, input contract.TokenPayloadInput) (tokenString string, resp contract.TokenPayload, err error) {
	payload := newPayload(fromContractTokenPayloadInput(input), accessTokenDurationTime)

	tokenString, err = p.createToken(ctx, payload)
	if err != nil {
		return tokenString, resp, err
	}

	return tokenString, payload.toContract(), nil
}

func (p *pasetoAuth) CreateRefreshToken(ctx context.Context, input contract.TokenPayloadInput) (tokenString string, resp contract.TokenPayload, err error) {
	payload := newPayload(fromContractTokenPayloadInput(input), refreshTokenDurationTime)

	tokenString, err = p.createToken(ctx, payload)
	if err != nil {
		return tokenString, resp, err
	}

	return tokenString, payload.toContract(), nil
}

func (p *pasetoAuth) VerifyToken(ctx context.Context, tokenStr string) (resp contract.TokenPayload, err error) {
	if strings.TrimSpace(tokenStr) == "" {
		return resp, apperr.ErrTokenInvalid
	}

	token, err := p.parser.ParseV4Local(p.key, tokenStr, nil)
	if err != nil {
		p.log.Error(ctx, "error to decrypt token", logger.Err(err))
		return resp, apperr.ErrTokenInvalid
	}

	payload, err := tokenPayloadFromClaims(token)
	if err != nil {
		p.log.Error(ctx, "error to parse token claims", logger.Err(err))
		return resp, apperr.ErrTokenInvalid
	}

	if err := payload.Valid(); err != nil {
		if errors.Is(err, errExpiredToken) {
			p.log.Warn(ctx, "token has expired")
			return resp, apperr.ErrTokenExpired
		}
		p.log.Error(ctx, "error to validate token", logger.Err(err))
		return resp, apperr.ErrTokenInvalid
	}

	return payload.toContract(), nil
}

func (p *pasetoAuth) createToken(_ context.Context, payload *tokenPayload) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	if err := token.Set(claimData, payload); err != nil {
		return "", err
	}

	return token.V4Encrypt(p.key, nil), nil
}

func tokenPayloadFromClaims(token *paseto.Token) (*tokenPayload, error) {
	payload := &tokenPayload{}
	if err := token.Get(claimData, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
