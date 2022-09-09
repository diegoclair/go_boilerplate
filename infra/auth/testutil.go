package auth

import (
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/stretchr/testify/require"
)

type utilArgs struct {
	accountUUID          string
	sessionUUID          string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	tokenType            string
	expiredToken         bool
	withoutPrivateKey    bool
	wantErr              bool
	wantErrValue         error
}

func getConfig(t *testing.T, args utilArgs) config.Config {
	cfgPointer, err := config.GetConfigEnvironment("../../" + config.ConfigDefaultName)
	cfg := *cfgPointer //do not use pointer here because GetConfigEnvironment return the same pointer and it can generate problem when run multiple tests
	cfg.App.Auth.AccessTokenType = args.tokenType

	if args.accessTokenDuration.String() != "0s" {
		cfg.App.Auth.AccessTokenDuration = args.accessTokenDuration
	}
	if args.refreshTokenDuration.String() != "0s" {
		cfg.App.Auth.RefreshTokenDuration = args.refreshTokenDuration
	}

	if args.expiredToken {
		cfg.App.Auth.AccessTokenDuration = 0 * time.Second
		cfg.App.Auth.RefreshTokenDuration = 0 * time.Second
	}

	if args.withoutPrivateKey {
		cfg.App.Auth.JWTPrivateKey = ""
		cfg.App.Auth.PasetoSymmetricKey = ""
	}
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg)
	return cfg
}

func createTestAccessToken(t *testing.T, maker AuthToken, args utilArgs) (token string, tokenPayload *tokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateAccessToken(args.accountUUID, args.sessionUUID)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func createTestRefreshToken(t *testing.T, maker AuthToken, args utilArgs) (token string, tokenPayload *tokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateRefreshToken(args.accountUUID, args.sessionUUID)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func validateTokenCreation(t *testing.T, args utilArgs, token string, tokenPayload *tokenPayload, err error) {
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)
	require.NotEmpty(t, tokenPayload)

	require.Equal(t, args.accountUUID, tokenPayload.AccountUUID)
	require.Equal(t, args.sessionUUID, tokenPayload.SessionUUID)
	require.NotZero(t, tokenPayload.IssuedAt)
	require.NotZero(t, tokenPayload.ExpiredAt)
}

func validateTokenMaker(t *testing.T, args utilArgs) {

	cfg := getConfig(t, args)
	require.NotNil(t, cfg)

	maker, err := NewAuthToken(cfg.App.Auth)
	if args.wantErr != (err != nil) {
		t.Errorf("NewAuthToken() error = %v, wantErr %v", err, args.wantErr)
	}
	require.Equal(t, args.wantErrValue, err)

	if args.wantErr {
		return
	}
	createTestAccessToken(t, maker, args)
	createTestRefreshToken(t, maker, args)
}
