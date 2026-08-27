package auth

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testNamespace   = "account"
	testMaxAttempts = 5
	testWindow      = time.Hour
)

func newTestLockout(t *testing.T, ctrl *gomock.Controller) (*Lockout, *mocks.MockLoginAttemptStore) {
	t.Helper()

	store := mocks.NewMockLoginAttemptStore(ctrl)
	lockout, err := NewLockout(store, testNamespace, testMaxAttempts, testWindow)
	require.NoError(t, err)

	return lockout, store
}

// A dependency accepted as nil trades a refused boot for a defence that is
// silently off in production.
func TestNewLockout_RefusesAnUnusableConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockLoginAttemptStore(ctrl)

	tests := []struct {
		name        string
		store       contract.LoginAttemptStore
		namespace   string
		maxAttempts int
		window      time.Duration
	}{
		{"no store", nil, testNamespace, testMaxAttempts, testWindow},
		{"blank namespace", store, "   ", testMaxAttempts, testWindow},
		{"no attempts", store, testNamespace, 0, testWindow},
		{"no window", store, testNamespace, testMaxAttempts, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLockout(tt.store, tt.namespace, tt.maxAttempts, tt.window)

			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

// Case and padding are the cheapest way to buy a fresh set of attempts, so every
// spelling of one identity has to land on one key.
func TestLockout_EverySpellingOfAnIdentitySharesOneKey(t *testing.T) {
	tests := []struct {
		name      string
		canonical string
		spellings []string
	}{
		{
			name:      "padding",
			canonical: "01234567890",
			spellings: []string{" 01234567890 ", "01234567890", "\t01234567890\n"},
		},
		{
			name:      "case",
			canonical: "katia@x.com",
			spellings: []string{"Katia@X.com", "KATIA@X.COM", " katia@x.com "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lockout, store := newTestLockout(t, ctrl)

			var keys []string
			capture := func(_ context.Context, key string) (int64, error) {
				keys = append(keys, key)
				return 0, contract.ErrCacheMiss
			}

			for _, spelling := range tt.spellings {
				store.EXPECT().GetInt(gomock.Any(), gomock.Any()).DoAndReturn(capture)
				_, err := lockout.Check(context.Background(), spelling)
				require.NoError(t, err)
			}

			require.Len(t, keys, len(tt.spellings))
			for _, key := range keys {
				assert.Equal(t, testNamespace+":login_attempt:"+tt.canonical, key)
			}
		})
	}
}

func TestLockout_Check(t *testing.T) {
	t.Run("no counter yet leaves every attempt available", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lockout, store := newTestLockout(t, ctrl)

		store.EXPECT().GetInt(gomock.Any(), gomock.Any()).Return(int64(0), contract.ErrCacheMiss)

		status, err := lockout.Check(context.Background(), "01234567890")

		require.NoError(t, err)
		assert.False(t, status.Locked)
		assert.Equal(t, testMaxAttempts, status.Remaining)
	})

	// The ceiling counts failures already made, so the identity still owns the
	// try that reaches it — shifting this by one silently costs a live attempt.
	t.Run("one failure below the ceiling is still open", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lockout, store := newTestLockout(t, ctrl)

		store.EXPECT().GetInt(gomock.Any(), gomock.Any()).Return(int64(testMaxAttempts-1), nil)

		status, err := lockout.Check(context.Background(), "01234567890")

		require.NoError(t, err)
		assert.False(t, status.Locked)
		assert.Equal(t, 1, status.Remaining)
	})

	t.Run("at the ceiling it is locked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lockout, store := newTestLockout(t, ctrl)

		store.EXPECT().GetInt(gomock.Any(), gomock.Any()).Return(int64(testMaxAttempts), nil)

		status, err := lockout.Check(context.Background(), "01234567890")

		require.NoError(t, err)
		assert.True(t, status.Locked)
	})

	t.Run("a store failure is the caller's to decide on", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		lockout, store := newTestLockout(t, ctrl)

		store.EXPECT().GetInt(gomock.Any(), gomock.Any()).Return(int64(0), assert.AnError)

		status, err := lockout.Check(context.Background(), "01234567890")

		require.ErrorIs(t, err, assert.AnError)
		assert.False(t, status.Locked)
	})
}

// Without the window the key never expires, and an identity that has no account
// is locked out of opening one.
func TestLockout_FailCountsWithinTheWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockout, store := newTestLockout(t, ctrl)

	var ttl time.Duration
	store.EXPECT().IncrBy(gomock.Any(), testNamespace+":login_attempt:01234567890", int64(1), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, _ int64, expiration time.Duration) (int64, error) {
			ttl = expiration
			return 3, nil
		})

	status, err := lockout.Fail(context.Background(), " 01234567890 ")

	require.NoError(t, err)
	assert.Equal(t, testWindow, ttl)
	assert.False(t, status.Locked)
	assert.Equal(t, testMaxAttempts-3, status.Remaining)
}

func TestLockout_FailLocksOnTheCountThatReachesTheCeiling(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockout, store := newTestLockout(t, ctrl)

	store.EXPECT().IncrBy(gomock.Any(), gomock.Any(), int64(1), testWindow).
		Return(int64(testMaxAttempts), nil)

	status, err := lockout.Fail(context.Background(), "01234567890")

	require.NoError(t, err)
	assert.True(t, status.Locked)
	assert.Zero(t, status.Remaining)
}

func TestLockout_FailKeepsTheStoreFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockout, store := newTestLockout(t, ctrl)

	store.EXPECT().IncrBy(gomock.Any(), gomock.Any(), int64(1), testWindow).
		Return(int64(0), assert.AnError)

	status, err := lockout.Fail(context.Background(), "01234567890")

	require.ErrorIs(t, err, assert.AnError)
	assert.False(t, status.Locked)
	assert.Zero(t, status.Remaining)
}

func TestLockout_ResetClearsTheCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockout, store := newTestLockout(t, ctrl)

	store.EXPECT().Delete(gomock.Any(), testNamespace+":login_attempt:01234567890").Return(nil)

	require.NoError(t, lockout.Reset(context.Background(), " 01234567890 "))
}
