package cache

type errorString string

func (err errorString) Error() string {
	return string(err)
}

// ErrCache represents cache manager related errors
const (
	// ErrCacheMiss indicates a cache miss when fetching an item from CacheManager.
	ErrCacheMiss = errorString("cache miss: key not found")
)
