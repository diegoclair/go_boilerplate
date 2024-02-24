package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"

	"github.com/redis/go-redis/v9"
)

// IRedisCache is the interface for the RedisCache - it's used to help testing
type IRedisCache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

// redisCache implements the CacheManager interface
type redisCache struct {
	cfg   config.RedisConfig
	redis IRedisCache
	log   logger.Logger
}

// newRedisCache returns a RedisCache instance
func newRedisCache(ctx context.Context, cfg *config.Config, log logger.Logger) (*redisCache, error) {
	rCfg := cfg.Cache.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     rCfg.Host + ":" + strconv.Itoa(rCfg.Port),
		Password: rCfg.Pass,
		DB:       rCfg.DB,
	})

	cfg.AddCloser(func() {
		log.Info(ctx, "Closing redis connection...")
		if err := client.Close(); err != nil {
			log.Errorf(ctx, "Error closing redis connection: %v", err)
		}
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &redisCache{
		cfg:   rCfg,
		redis: client,
		log:   log,
	}, nil
}

// GetItem returns an Item from cache
func (r *redisCache) GetItem(ctx context.Context, key string) (data []byte, err error) {
	val, err := r.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return val, ErrCacheMiss
	} else if err != nil {
		return val, err
	}

	return val, nil
}

// SetItem sets an item in cache
func (r *redisCache) SetItem(ctx context.Context, key string, data []byte) error {
	err := r.SetItemWithExpiration(ctx, key, data, r.cfg.DefaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetItemWithExpiration sets an item in cache
func (r *redisCache) SetItemWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	err := r.redis.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetInt returns an int64 from cache
func (r *redisCache) GetInt(ctx context.Context, key string) (data int64, err error) {
	val, err := r.GetItem(ctx, key)
	if err == ErrCacheMiss {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return strconv.ParseInt(string(val), 10, 64)
}

// GetString returns an string from cache
func (r *redisCache) GetString(ctx context.Context, key string) (data string, err error) {
	val, err := r.GetItem(ctx, key)
	if err == ErrCacheMiss {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return string(val), nil
}

// SetString sets an item in cache
func (r *redisCache) SetString(ctx context.Context, key string, data string) error {
	err := r.SetStringWithExpiration(ctx, key, data, r.cfg.DefaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStringWithExpiration sets an item in cache
func (r *redisCache) SetStringWithExpiration(ctx context.Context, key string, data string, expiration time.Duration) error {
	err := r.SetItemWithExpiration(ctx, key, []byte(data), expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetStruct returns an struct from cache
func (r *redisCache) GetStruct(ctx context.Context, key string, data interface{}) (err error) {
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

// SetStruct sets an item in cache
func (r *redisCache) SetStruct(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = r.cfg.DefaultExpiration
	}
	err := r.SetStructWithExpiration(ctx, key, data, expiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStructWithExpiration sets an item in cache
func (r *redisCache) SetStructWithExpiration(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
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
func (r *redisCache) GetExpiration(ctx context.Context, key string) (expiration time.Duration, err error) {
	expiration, err = r.redis.TTL(ctx, key).Result()
	if err != nil {
		return expiration, err
	}

	return expiration, nil
}

// SetExpiration sets the expiration time for a key
func (r *redisCache) SetExpiration(ctx context.Context, key string, expiration time.Duration) (err error) {
	err = r.redis.Expire(ctx, key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// Increase increases an int key, setting it to zero if the key doesn't exists
func (r *redisCache) Increase(ctx context.Context, key string) (err error) {
	err = r.redis.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a list of keys from the cache
func (r *redisCache) Delete(ctx context.Context, keys ...string) (err error) {
	err = r.redis.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

// CleanAll clean everything with the prefix.
func (r *redisCache) CleanAll(ctx context.Context) (err error) {
	keys, err := r.redis.Keys(ctx, "*").Result()
	if err == redis.Nil {
		return ErrCacheMiss
	} else if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = r.redis.Del(ctx, keys...).Err()
		if err == redis.Nil {
			return ErrCacheMiss
		} else if err != nil {
			return err
		}
	}

	return nil
}
