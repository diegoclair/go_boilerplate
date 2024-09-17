package auth

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/stretchr/testify/require"
)

type utilArgs struct {
	emptyToken           bool
	payload              contract.TokenPayloadInput
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	expiredToken         bool
	withoutPrivateKey    bool
	wantErr              bool
	wantErrValue         error
}

func getConfig(t *testing.T, args utilArgs) *configmock.ConfigMock {
	cfg := configmock.New()

	if args.accessTokenDuration.String() != "0s" {
		cfg.Auth.AccessTokenDuration = args.accessTokenDuration
	}
	if args.refreshTokenDuration.String() != "0s" {
		cfg.Auth.RefreshTokenDuration = args.refreshTokenDuration
	}

	if args.expiredToken {
		cfg.Auth.AccessTokenDuration = 0 * time.Second
		cfg.Auth.RefreshTokenDuration = 0 * time.Second
	}

	if args.withoutPrivateKey {
		cfg.Auth.PasetoSymmetricKey = ""
	}

	require.NotEmpty(t, cfg)
	return cfg
}

func createTestAccessToken(ctx context.Context, t *testing.T, maker contract.AuthToken, args utilArgs) (token string, tokenPayload contract.TokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateAccessToken(ctx, args.payload)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func createTestRefreshToken(ctx context.Context, t *testing.T, maker contract.AuthToken, args utilArgs) (token string, tokenPayload contract.TokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateRefreshToken(ctx, args.payload)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func validateTokenCreation(t *testing.T, args utilArgs, token string, tokenPayload contract.TokenPayload, err error) {
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)
	require.NotEmpty(t, tokenPayload)

	require.Equal(t, args.payload.AccountUUID, tokenPayload.AccountUUID)
	require.Equal(t, args.payload.SessionUUID, tokenPayload.SessionUUID)
	require.NotZero(t, tokenPayload.IssuedAt)
	require.NotZero(t, tokenPayload.ExpiredAt)
}

func validateTokenMaker(t *testing.T, args utilArgs) {

	ctx := context.Background()
	cfg := getConfig(t, args)
	require.NotNil(t, cfg)

	maker, err := getTokenAuth(cfg)
	if args.wantErr {
		require.Equal(t, args.wantErrValue, err)
		return
	}
	require.NoError(t, err)

	createTestAccessToken(ctx, t, maker, args)
	createTestRefreshToken(ctx, t, maker, args)
}

func getTokenAuth(cfg *configmock.ConfigMock) (contract.AuthToken, error) {
	return NewAuthToken(cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
		cfg.Auth.PasetoSymmetricKey,
		cfg.GetLogger(),
	)
}
