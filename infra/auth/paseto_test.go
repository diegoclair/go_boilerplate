package auth

import (
	"testing"
	"time"
)

func TestPasetoTokenMaker(t *testing.T) {

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should create token without error",
			args: args{
				accountUUID:          "account-123",
				sessionUUID:          "session-123",
				accessTokenDuration:  time.Second,
				refreshTokenDuration: time.Second * 2,
			},
		},
		{
			name: "Should return error with an invalid privateKey",
			args: args{
				withoutPrivateKey: true,
				wantErr:           true,
				wantErrValue:      errInvalidPrivateKey,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateTokenMaker(t, tt.args, tokenTypePaseto)
		})
	}
}
