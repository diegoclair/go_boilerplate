package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/diegoclair/go_utils/logger"

	"github.com/redis/go-redis/v9"
)

// IRedisCache is the interface for the RedisCache - it's used to help testing
type IRedisCache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

// CacheManager implements the CacheManager interface
type CacheManager struct {
	defaultExpiration time.Duration
	redis             IRedisCache
	log               logger.Logger
}

// NewRedisCache returns a RedisCache instance
func NewRedisCache(ctx context.Context, addr string, password string, db int, defaultExpiration time.Duration, log logger.Logger) (*CacheManager, *redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, nil, err
	}

	return &CacheManager{
		redis:             client,
		defaultExpiration: defaultExpiration,
		log:               log,
	}, client, nil
}

// Set stores a value in cache. Accepts string, []byte, int, int64 or any struct (JSON marshaled).
// Expiration is optional: omit for default, pass for custom.
func (r *CacheManager) Set(ctx context.Context, key string, data any, expiration ...time.Duration) error {
	exp := r.defaultExpiration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	var bytes []byte
	switch v := data.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case int:
		bytes = []byte(strconv.Itoa(v))
	case int64:
		bytes = []byte(strconv.FormatInt(v, 10))
	default:
		var err error
		bytes, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	return r.redis.Set(ctx, key, bytes, exp).Err()
}

// Get returns raw bytes from cache
func (r *CacheManager) Get(ctx context.Context, key string) (data []byte, err error) {
	val, err := r.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return val, ErrCacheMiss
	} else if err != nil {
		return val, err
	}

	return val, nil
}

// GetInt returns an int64 from cache
func (r *CacheManager) GetInt(ctx context.Context, key string) (data int64, err error) {
	val, err := r.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return strconv.ParseInt(string(val), 10, 64)
}

// GetString returns a string from cache
func (r *CacheManager) GetString(ctx context.Context, key string) (data string, err error) {
	val, err := r.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return string(val), nil
}

// GetStruct receives a pointer to a struct and populates it from cache (JSON unmarshal)
func (r *CacheManager) GetStruct(ctx context.Context, key string, data any) (err error) {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(val, &data)
	if err != nil {
		return err
	}

	return nil
}

// GetExpiration returns the expiration time for a key
func (r *CacheManager) GetExpiration(ctx context.Context, key string) (expiration time.Duration, err error) {
	expiration, err = r.redis.TTL(ctx, key).Result()
	if err != nil {
		return expiration, err
	}

	return expiration, nil
}

// SetExpiration sets the expiration time for a key
func (r *CacheManager) SetExpiration(ctx context.Context, key string, expiration time.Duration) (err error) {
	err = r.redis.Expire(ctx, key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// Increase increases an int key, setting it to zero if the key doesn't exists
func (r *CacheManager) Increase(ctx context.Context, key string) (err error) {
	err = r.redis.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a list of keys from the cache
func (r *CacheManager) Delete(ctx context.Context, keys ...string) (err error) {
	err = r.redis.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

// CleanAll clean everything with the prefix.
func (r *CacheManager) CleanAll(ctx context.Context) (err error) {
	keys, err := r.redis.Keys(ctx, "*").Result()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	} else if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = r.redis.Del(ctx, keys...).Err()
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		} else if err != nil {
			return err
		}
	}

	return nil
}

// GetAllKeys retrieves all keys matching a given pattern
func (r *CacheManager) GetAllKeys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	return keys, nil
}
