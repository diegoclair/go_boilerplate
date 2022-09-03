package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/stretchr/testify/require"
)

type args struct {
	accountUUID          string
	sessionUUID          string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	withoutPrivateKey    bool
	wantErr              bool
	wantErrValue         error
}

func getConfig(t *testing.T) *config.Config {
	cfg, err := config.GetConfigEnvironment("../../" + config.ConfigDefaultName)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg)
	return cfg
}

func createAccessToken(t *testing.T, maker AuthToken, args args) (token string, tokenPayload *tokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateAccessToken(args.accountUUID, args.sessionUUID)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func createRefreshToken(t *testing.T, maker AuthToken, args args) (token string, tokenPayload *tokenPayload) {
	var err error
	token, tokenPayload, err = maker.CreateRefreshToken(args.accountUUID, args.sessionUUID)
	validateTokenCreation(t, args, token, tokenPayload, err)
	return
}

func validateTokenCreation(t *testing.T, args args, token string, tokenPayload *tokenPayload, err error) {
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, tokenPayload)
	require.NotEmpty(t, tokenPayload)

	require.Equal(t, args.accountUUID, tokenPayload.AccountUUID)
	require.Equal(t, args.sessionUUID, tokenPayload.SessionUUID)
	require.NotZero(t, tokenPayload.IssuedAt)
	require.NotZero(t, tokenPayload.ExpiredAt)
}

func validateTokenMaker(t *testing.T, args args, tokenType string) {

	cfg := getConfig(t)
	require.NotNil(t, cfg)

	configAuth := config.AuthConfig{
		AccessTokenType:      tokenType,
		AccessTokenDuration:  args.accessTokenDuration,
		RefreshTokenDuration: args.refreshTokenDuration,
	}
	if !args.withoutPrivateKey {
		if tokenType == tokenTypeJWT {
			configAuth.JWTPrivateKey = cfg.App.Auth.JWTPrivateKey
		}
		if tokenType == tokenTypePaseto {
			configAuth.PasetoSymmetricKey = cfg.App.Auth.PasetoSymmetricKey
		}
	}

	maker, err := NewAuthToken(configAuth)
	if args.wantErr != (err != nil) {
		fmt.Println("diegooo ", cfg.App.Auth.PasetoSymmetricKey)
		t.Errorf("NewAuthToken() error = %v, wantErr %v", err, args.wantErr)
	}
	require.Equal(t, args.wantErrValue, err)

	if args.wantErr {
		return
	}
	createAccessToken(t, maker, args)
	createRefreshToken(t, maker, args)
}
