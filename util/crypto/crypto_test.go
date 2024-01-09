package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	c := NewCrypto()
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should return a hashed password",
			args:    args{password: "123456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEmpty(t, got)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	c := NewCrypto()

	type args struct {
		password       string
		hashedPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should return nil if the password is correct",
			args:    args{password: "123456", hashedPassword: "$2a$10$jN9Oi/xk63jSMlWFHvqSseIcJh/5YNfGfpkd9VpKndvQzhJChYUAW"},
			wantErr: false,
		},
		{
			name:    "Should return an error if the hashed password is incorrect",
			args:    args{password: "123456", hashedPassword: "$2a$10$2H2yQjIe2m4n1Yh1uV4f3u3Z4K6d1Qa1c1f2v3e4r5t6y7u8i9o0"},
			wantErr: true,
		},
		{
			name:    "Should return an error if the hashed password is incorrect",
			args:    args{password: "123456", hashedPassword: "$2a$10$2H2yQjIe2m4n1Yh1uV4f3u3Z4K6d1Qa1c1f2v3e4r5t6y7u8i9o0p"},
			wantErr: true,
		},
		{
			name:    "Should return an error if the hashed password is incorrect",
			args:    args{password: "123456", hashedPassword: "$2a$10$2H2yQjIe2m4n1Yh1uV4f3u3Z4K6d1Qa1c1f2v3e4r5t6y7u8i9o0p"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CheckPassword(tt.args.password, tt.args.hashedPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
