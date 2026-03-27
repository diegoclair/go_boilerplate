package logger

import (
	"context"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test", true)
	require.NotNil(t, logger)
}

func TestGetContextValue(t *testing.T) {
	t.Run("Should return zero value and false when context is nil", func(t *testing.T) {
		var ctx context.Context = nil
		value, ok := getContextValue[string](ctx, infra.SessionKey)
		require.False(t, ok)
		require.Equal(t, "", value)
	})

	t.Run("Should return zero value and false when key not in context", func(t *testing.T) {
		ctx := context.Background()
		value, ok := getContextValue[string](ctx, infra.SessionKey)
		require.False(t, ok)
		require.Equal(t, "", value)
	})

	t.Run("Should return string value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")
		value, ok := getContextValue[string](ctx, infra.SessionKey)
		require.True(t, ok)
		require.Equal(t, "sessionCode", value)
	})

	t.Run("Should return false when type does not match", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, 123)
		value, ok := getContextValue[string](ctx, infra.SessionKey)
		require.False(t, ok)
		require.Equal(t, "", value)
	})
}

func TestAddDefaultAttributesToLogger(t *testing.T) {
	t.Run("Should return session and account_uuid attributes", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")
		ctx = context.WithValue(ctx, infra.AccountUUIDKey, "accountUUID")

		args := addDefaultAttributesToLogger(ctx)
		require.Len(t, args, 2)
	})

	t.Run("Should return only session attribute when account_uuid is missing", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.SessionKey, "sessionCode")

		args := addDefaultAttributesToLogger(ctx)
		require.Len(t, args, 1)
	})

	t.Run("Should return only account_uuid attribute when session is missing", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), infra.AccountUUIDKey, "accountUUID")

		args := addDefaultAttributesToLogger(ctx)
		require.Len(t, args, 1)
	})

	t.Run("Should return empty when context has no values", func(t *testing.T) {
		ctx := context.Background()
		args := addDefaultAttributesToLogger(ctx)
		require.Empty(t, args)
	})
}
