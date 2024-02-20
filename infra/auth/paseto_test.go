package auth

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/stretchr/testify/require"
)

func TestPasetoTokenMaker(t *testing.T) {

	tests := []struct {
		name string
		args utilArgs
	}{
		{
			name: "Should create token without error",
			args: utilArgs{
				accountUUID:          "account-123",
				sessionUUID:          "session-123",
				accessTokenDuration:  time.Second,
				refreshTokenDuration: time.Second * 2,
			},
		},
		{
			name: "Should return error with an invalid privateKey",
			args: utilArgs{
				withoutPrivateKey: true,
				wantErr:           true,
				wantErrValue:      errInvalidPrivateKey,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.tokenType = tokenTypePaseto
			validateTokenMaker(t, tt.args)
		})
	}
}

func Test_paseto_VerifyToken(t *testing.T) {

	tests := []struct {
		name    string
		args    utilArgs
		wantErr bool
	}{
		{
			name: "Should pass without error",
		},
		{
			name: "Should return error for a expired token",
			args: utilArgs{
				expiredToken: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.tokenType = tokenTypePaseto
			cfg := getConfig(t, tt.args)
			maker, err := NewAuthToken(cfg.App.Auth, logger.NewNoop())
			require.NoError(t, err)

			ctx := context.Background()

			token, tokenPayload := createTestAccessToken(ctx, t, maker, tt.args)

			gotPayload, err := maker.VerifyToken(ctx, token)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtAuth.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.args.sessionUUID, gotPayload.SessionUUID)
			require.Equal(t, tt.args.accountUUID, gotPayload.AccountUUID)
			require.WithinDuration(t, tokenPayload.IssuedAt, gotPayload.IssuedAt, 1*time.Second)
			require.WithinDuration(t, tokenPayload.ExpiredAt, gotPayload.ExpiredAt, 1*time.Second)
		})
	}
}
