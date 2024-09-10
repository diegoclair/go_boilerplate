package auth

import (
	"context"
	"testing"
	"time"

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
				payload: TokenPayloadInput{
					AccountUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
					SessionUUID: "d152a340-9a87-4d32-85ad-19df4c9934cd",
				},
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
			validateTokenMaker(t, tt.args)
		})
	}
}

func Test_paseto_VerifyToken(t *testing.T) {
	tests := []struct {
		name string
		args utilArgs
	}{
		{
			name: "Should pass without error",
		},
		{
			name: "Should return error for a expired token",
			args: utilArgs{
				expiredToken: true,
				wantErr:      true,
			},
		},
		{
			name: "Should return error for a empty token",
			args: utilArgs{
				emptyToken:   true,
				wantErr:      true,
				wantErrValue: errInvalidToken,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maker, err := getTokenAuth(getConfig(t, tt.args))
			require.NoError(t, err)

			ctx := context.Background()

			token, tokenPayload := createTestAccessToken(ctx, t, maker, tt.args)
			if tt.args.emptyToken {
				token = ""
			}

			gotPayload, err := maker.VerifyToken(ctx, token)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("jwtAuth.VerifyToken() error = %v, wantErr %v", err, tt.args.wantErr)
				return
			}
			if tt.args.wantErr {
				return
			}
			require.Equal(t, tt.args.payload.SessionUUID, gotPayload.SessionUUID)
			require.Equal(t, tt.args.payload.AccountUUID, gotPayload.AccountUUID)
			require.WithinDuration(t, tokenPayload.IssuedAt, gotPayload.IssuedAt, 1*time.Second)
			require.WithinDuration(t, tokenPayload.ExpiredAt, gotPayload.ExpiredAt, 1*time.Second)
		})
	}
}
