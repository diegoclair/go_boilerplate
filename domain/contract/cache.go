package contract

import "time"

// CacheManager defines the main caching interface
//   - Get methods can return domain.ErrCacheMiss
type CacheManager interface {
	GetItem(key string) ([]byte, error)
	SetItem(key string, data []byte) error
	SetItemWithExpiration(key string, data []byte, expiration time.Duration) error

	GetString(key string) (string, error)
	SetString(key string, data string) error
	SetStringWithExpiration(key string, data string, expiration time.Duration) error

	GetInt(key string) (data int64, err error)
	Increase(key string) error

	GetStruct(key string, data interface{}) error
	SetStruct(key string, data interface{}, expiration time.Duration) error
	SetStructWithExpiration(key string, data interface{}, expiration time.Duration) error

	GetExpiration(key string) (time.Duration, error)
	SetExpiration(key string, expiration time.Duration) error

	Delete(keys ...string) error
	CleanAll() error
}
