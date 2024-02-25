package auth

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
	"github.com/stretchr/testify/require"
)

type utilArgs struct {
	emptyToken           bool
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

func getConfig(t *testing.T, args utilArgs) *config.Config {
	cfgPointer, err := config.GetConfigEnvironment(config.ProfileTest)
	require.NoError(t, err)
	require.NotNil(t, cfgPointer)

	cfg := copyConfig(cfgPointer)
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

	require.NotEmpty(t, cfg)
	return cfg
}

// copyConfig returns a copy of the config avoiding problems with pointer and multiple tests
func copyConfig(cfg *config.Config) *config.Config {
	return &config.Config{
		App:   cfg.App,
		Cache: cfg.Cache,
		DB:    cfg.DB,
		Log:   cfg.Log,
	}
}

func createTestAccessToken(ctx context.Context, t *testing.T, maker AuthToken, args utilArgs) (token string, tokenPayload *TokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateAccessToken(ctx, TokenPayloadInput{AccountUUID: args.accountUUID, SessionUUID: args.sessionUUID})
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func createTestRefreshToken(ctx context.Context, t *testing.T, maker AuthToken, args utilArgs) (token string, tokenPayload *TokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateRefreshToken(ctx, TokenPayloadInput{AccountUUID: args.accountUUID, SessionUUID: args.sessionUUID})
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func validateTokenCreation(t *testing.T, args utilArgs, token string, tokenPayload *TokenPayload, err error) {
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

	ctx := context.Background()
	cfg := getConfig(t, args)
	require.NotNil(t, cfg)

	maker, err := NewAuthToken(cfg.App.Auth, logger.NewNoop())
	if args.wantErr != (err != nil) {
		t.Errorf("NewAuthToken() error = %v, wantErr %v", err, args.wantErr)
	}
	require.Equal(t, args.wantErrValue, err)

	if args.wantErr {
		return
	}
	createTestAccessToken(ctx, t, maker, args)
	createTestRefreshToken(ctx, t, maker, args)
}
