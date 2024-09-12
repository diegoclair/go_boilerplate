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

// GetItem returns an Item from cache
func (r *CacheManager) GetItem(ctx context.Context, key string) (data []byte, err error) {
	val, err := r.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return val, ErrCacheMiss
	} else if err != nil {
		return val, err
	}

	return val, nil
}

// SetItem sets an item in cache
func (r *CacheManager) SetItem(ctx context.Context, key string, data []byte) error {
	err := r.SetItemWithExpiration(ctx, key, data, r.defaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetItemWithExpiration sets an item in cache
func (r *CacheManager) SetItemWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	err := r.redis.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetInt returns an int64 from cache
func (r *CacheManager) GetInt(ctx context.Context, key string) (data int64, err error) {
	val, err := r.GetItem(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return strconv.ParseInt(string(val), 10, 64)
}

// GetString returns an string from cache
func (r *CacheManager) GetString(ctx context.Context, key string) (data string, err error) {
	val, err := r.GetItem(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return string(val), nil
}

// SetString sets an item in cache with default expiration
func (r *CacheManager) SetString(ctx context.Context, key string, data string) error {
	err := r.SetStringWithExpiration(ctx, key, data, r.defaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStringWithExpiration sets an item in cache with an expiration
func (r *CacheManager) SetStringWithExpiration(ctx context.Context, key string, data string, expiration time.Duration) error {
	err := r.SetItemWithExpiration(ctx, key, []byte(data), expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetStruct receive a pointer to a struct and returns the struct from cache
func (r *CacheManager) GetStruct(ctx context.Context, key string, data any) (err error) {
	val, err := r.GetItem(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(val, &data)
	if err != nil {
		return err
	}

	return nil
}

// SetStruct sets an item in cache with default expiration
func (r *CacheManager) SetStruct(ctx context.Context, key string, data any) error {
	err := r.SetStructWithExpiration(ctx, key, data, r.defaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStructWithExpiration sets an item in cache
func (r *CacheManager) SetStructWithExpiration(ctx context.Context, key string, data any, expiration time.Duration) error {
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.SetItemWithExpiration(ctx, key, dataString, expiration)
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
