package contract

import (
	"context"
	"time"
)

type errorString string

func (err errorString) Error() string {
	return string(err)
}

// Lives next to the ports so a consumer can tell an absent key from a broken
// store without depending on the cache implementation.
const ErrCacheMiss = errorString("cache miss: key not found")

// The slice of the cache the login lockout needs, so the shared sign-in flow
// does not depend on the whole cache manager.
type LoginAttemptStore interface {
	GetInt(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, delta int64, ttl time.Duration) (int64, error)
	Delete(ctx context.Context, keys ...string) error
}
