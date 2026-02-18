package contract

import (
	"context"
	"time"
)

// CacheManager defines the main caching interface
//   - Get methods can return cache.ErrCacheMiss
//   - Set accepts string, []byte, int, int64 or struct (JSON marshaled internally)
//   - Set expiration is variadic: omit for default, pass for custom
type CacheManager interface {
	Set(ctx context.Context, key string, data any, expiration ...time.Duration) error

	Get(ctx context.Context, key string) ([]byte, error)
	GetString(ctx context.Context, key string) (string, error)
	GetInt(ctx context.Context, key string) (data int64, err error)
	GetStruct(ctx context.Context, key string, data any) error

	Increase(ctx context.Context, key string) error

	GetExpiration(ctx context.Context, key string) (time.Duration, error)
	SetExpiration(ctx context.Context, key string, expiration time.Duration) error

	Delete(ctx context.Context, keys ...string) error
	CleanAll(ctx context.Context) error
	GetAllKeys(ctx context.Context, pattern string) ([]string, error)
}
