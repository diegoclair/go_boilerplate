package auth

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/cache"
	"github.com/diegoclair/go_boilerplate/infra/configmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The window and the atomicity of the counter are Redis semantics, so a mock
// cannot prove either — only a real server can.
func TestLockout_AgainstRedis(t *testing.T) {
	ctx := context.Background()

	cfg := configmock.New()
	closeContainer := cache.SetRedisTestContainerConfig(ctx, cfg)
	t.Cleanup(closeContainer)

	redisCache, client, err := cache.NewRedisCache(ctx,
		cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.DefaultExpiration,
		cfg.GetLogger(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	newLockout := func(t *testing.T, window time.Duration) *Lockout {
		t.Helper()

		lockout, err := NewLockout(redisCache, testNamespace, testMaxAttempts, window)
		require.NoError(t, err)

		return lockout
	}

	t.Run("every spelling of an identity counts against one counter", func(t *testing.T) {
		lockout := newLockout(t, time.Minute)
		spellings := []string{" Katia@X.com ", "KATIA@X.COM", "katia@x.com", "Katia@x.com", "katia@X.com"}
		require.Len(t, spellings, testMaxAttempts)

		for i, spelling := range spellings {
			status, err := lockout.Fail(ctx, spelling)
			require.NoError(t, err)

			if i < testMaxAttempts-1 {
				assert.False(t, status.Locked)
				assert.Equal(t, testMaxAttempts-i-1, status.Remaining)
				continue
			}

			assert.True(t, status.Locked)
		}

		status, err := lockout.Check(ctx, "katia@x.com")
		require.NoError(t, err)
		assert.True(t, status.Locked)
	})

	// The expiry belongs to the increment that created the key: renewed on every
	// failure, the wait would never end for whoever keeps trying.
	t.Run("the counter expires on its own without being renewed", func(t *testing.T) {
		lockout := newLockout(t, 2*time.Second)

		for range testMaxAttempts {
			_, err := lockout.Fail(ctx, "expires@x.com")
			require.NoError(t, err)
		}

		status, err := lockout.Check(ctx, "expires@x.com")
		require.NoError(t, err)
		require.True(t, status.Locked)

		require.Eventually(t, func() bool {
			status, err := lockout.Check(ctx, "expires@x.com")
			return err == nil && !status.Locked
		}, 5*time.Second, 200*time.Millisecond)
	})

	// A reset that missed the key would send a reactivated account straight back
	// to the ceiling on its next correct password.
	t.Run("reset clears the counter", func(t *testing.T) {
		lockout := newLockout(t, time.Minute)

		for range testMaxAttempts {
			_, err := lockout.Fail(ctx, "reset@x.com")
			require.NoError(t, err)
		}

		require.NoError(t, lockout.Reset(ctx, " RESET@X.com "))

		status, err := lockout.Check(ctx, "reset@x.com")
		require.NoError(t, err)
		assert.False(t, status.Locked)
		assert.Equal(t, testMaxAttempts, status.Remaining)
	})
}
