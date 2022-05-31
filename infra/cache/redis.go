package cache

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/diegoclair/go-boilerplate/infra/logger"
	"github.com/diegoclair/go-boilerplate/util/config"

	redis "gopkg.in/redis.v5"
)

// RedisCache implements the CacheManager interface
type RedisCache struct {
	cfg   config.RedisConfig
	redis *redis.Client
	log   logger.Logger
}

// NewRedisCache returns a RedisCache instance
func NewRedisCache(cfg config.RedisConfig, log logger.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Pass,
		DB:       cfg.DB,
	})

	_, err := client.Ping().Result()

	return &RedisCache{
		cfg:   cfg,
		redis: client,
		log:   log,
	}, err
}

// GetItem returns an Item from cache
func (r *RedisCache) GetItem(key string) (data []byte, err error) {
	val, err := r.redis.Get(key).Bytes()
	if err == redis.Nil {
		return val, ErrCacheMiss
	} else if err != nil {
		return val, err
	}

	return val, nil
}

// SetItem sets an item in cache
func (r *RedisCache) SetItem(key string, data []byte) error {
	err := r.SetItemWithExpiration(key, data, r.cfg.DefaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetItemWithExpiration sets an item in cache
func (r *RedisCache) SetItemWithExpiration(key string, data []byte, expiration time.Duration) error {
	err := r.redis.Set(key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetInt returns an int64 from cache
func (r *RedisCache) GetInt(key string) (data int64, err error) {
	val, err := r.GetItem(key)
	if err == ErrCacheMiss {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return strconv.ParseInt(string(val), 10, 64)
}

// GetString returns an string from cache
func (r *RedisCache) GetString(key string) (data string, err error) {
	val, err := r.GetItem(key)
	if err == ErrCacheMiss {
		return data, ErrCacheMiss
	} else if err != nil {
		return data, err
	}

	return string(val), nil
}

// SetString sets an item in cache
func (r *RedisCache) SetString(key string, data string) error {
	err := r.SetStringWithExpiration(key, data, r.cfg.DefaultExpiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStringWithExpiration sets an item in cache
func (r *RedisCache) SetStringWithExpiration(key string, data string, expiration time.Duration) error {
	err := r.SetItemWithExpiration(key, []byte(data), expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetStruct returns an struct from cache
func (r *RedisCache) GetStruct(key string, data interface{}) (err error) {
	val, err := r.GetItem(key)
	if err == ErrCacheMiss {
		return ErrCacheMiss
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(val, &data)
	if err != nil {
		return err
	}

	return nil
}

// SetStruct sets an item in cache
func (r *RedisCache) SetStruct(key string, data interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = r.cfg.DefaultExpiration
	}
	err := r.SetStructWithExpiration(key, data, expiration)
	if err != nil {
		return err
	}

	return nil
}

// SetStructWithExpiration sets an item in cache
func (r *RedisCache) SetStructWithExpiration(key string, data interface{}, expiration time.Duration) error {
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.SetItemWithExpiration(key, dataString, expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetExpiration returns the expiration time for a key
func (r *RedisCache) GetExpiration(key string) (expiration time.Duration, err error) {
	expiration, err = r.redis.TTL(key).Result()
	if err != nil {
		return expiration, err
	}

	return expiration, nil
}

// SetExpiration sets the expiration time for a key
func (r *RedisCache) SetExpiration(key string, expiration time.Duration) (err error) {
	err = r.redis.Expire(key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// Increase increases an int key, setting it to zero if the key doesn't exists
func (r *RedisCache) Increase(key string) (err error) {
	err = r.redis.Incr(key).Err()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a list of keys from the cache
func (r *RedisCache) Delete(keys ...string) (err error) {
	redisKeys := make([]string, len(keys))
	for i, key := range keys {
		redisKeys[i] = key
	}
	err = r.redis.Del(redisKeys...).Err()
	if err != nil {
		return err
	}

	return nil
}

// CleanAll clean everything with the prefix.
func (r *RedisCache) CleanAll() (err error) {
	keys, err := r.redis.Keys("*").Result()
	if err == redis.Nil {
		return ErrCacheMiss
	} else if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = r.redis.Del(keys...).Err()
	}
	if err == redis.Nil {
		return ErrCacheMiss
	} else if err != nil {
		return err
	}

	return nil
}
