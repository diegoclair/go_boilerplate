package logger

import (
	"context"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	logger := New(cfg)
	require.NotNil(t, logger)
}

func TestAddDefaultAttributesToLogger(t *testing.T) {
	ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")
	ctx = context.WithValue(ctx, infra.AccountUUIDKey, "accountUUID")

	args := addDefaultAttributesToLogger(ctx)
	require.Equal(t, "session", args[0])
	require.Equal(t, "sessionCode", args[1])
	require.Equal(t, "account_uuid", args[2])
	require.Equal(t, "accountUUID", args[3])
}

func TestGetContextValue(t *testing.T) {
	t.Run("Should return empty string when context is nil", func(t *testing.T) {
		var ctx context.Context = nil
		value := getContextValue(ctx, infra.SessionKey)
		require.Equal(t, "", value)
	})

	t.Run("Should return empty string when value is nil", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, nil)
		value := getContextValue(ctx, infra.AccountUUIDKey)
		require.Equal(t, "", value)
	})

	t.Run("Should return value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")

		value := getContextValue(ctx, infra.SessionKey)
		require.Equal(t, "sessionCode", value)
	})
}

func TestGetSession(t *testing.T) {
	t.Run("Should return empty string when session is empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "")
		sessionCode, ok := getSession(ctx)
		require.Equal(t, "", sessionCode)
		require.False(t, ok)
	})

	t.Run("Should return session code", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")
		sessionCode, ok := getSession(ctx)
		require.Equal(t, "sessionCode", sessionCode)
		require.True(t, ok)
	})
}

func TestGetAccountUUID(t *testing.T) {
	t.Run("Should return empty string when accountUUID is empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.AccountUUIDKey, "")
		accountUUID, ok := getAccountUUID(ctx)
		require.Equal(t, "", accountUUID)
		require.False(t, ok)
	})

	t.Run("Should return accountUUID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.AccountUUIDKey, "accountUUID")
		accountUUID, ok := getAccountUUID(ctx)
		require.Equal(t, "accountUUID", accountUUID)
		require.True(t, ok)
	})
}
