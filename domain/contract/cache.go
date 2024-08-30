package contract

import (
	"context"
	"time"
)

// CacheManager defines the main caching interface
//   - Get methods can return domain.ErrCacheMiss
type CacheManager interface {
	GetItem(ctx context.Context, key string) ([]byte, error)
	SetItem(ctx context.Context, key string, data []byte) error
	SetItemWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error

	GetString(ctx context.Context, key string) (string, error)
	SetString(ctx context.Context, key string, data string) error
	SetStringWithExpiration(ctx context.Context, key string, data string, expiration time.Duration) error

	GetInt(ctx context.Context, key string) (data int64, err error)
	Increase(ctx context.Context, key string) error

	GetStruct(ctx context.Context, key string, data any) error
	SetStruct(ctx context.Context, key string, data any, expiration time.Duration) error
	SetStructWithExpiration(ctx context.Context, key string, data any, expiration time.Duration) error

	GetExpiration(ctx context.Context, key string) (time.Duration, error)
	SetExpiration(ctx context.Context, key string, expiration time.Duration) error

	Delete(ctx context.Context, keys ...string) error
	CleanAll(ctx context.Context) error
}
