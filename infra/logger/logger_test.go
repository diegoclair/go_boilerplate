package logger

import (
	"context"
	"testing"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_utils/logger"
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
		require.Equal(t, "session", args[0].(logger.StringField).Key)
		require.Equal(t, "sessionCode", args[0].(logger.StringField).Value)
		require.Equal(t, "account_uuid", args[1].(logger.StringField).Key)
		require.Equal(t, "accountUUID", args[1].(logger.StringField).Value)
	})

	t.Run("Should return empty when context has no values", func(t *testing.T) {
		ctx := context.Background()
		args := addDefaultAttributesToLogger(ctx)
		require.Empty(t, args)
	})
}
