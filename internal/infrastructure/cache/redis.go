package cache

import (
	"context"
	"encoding/json"
	"fmt"
	logger "sinibeli/internal/pkg/logging"
	"time"

	"sinibeli/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type CacheConfig struct {
	Addr     string
	Password string
	DB       int
}

const (
	UserFileListKey = "user:files:%s"
	FileMetadataKey = "file:metadata:%s"
	FileExistsKey   = "file:exists:%s"
	ProductListKey  = "products:list:%s"
	ProductKey      = "product:%s"
	UserProfileKey  = "user:profile:%s"
)

const (
	FileMetadataTTL = 1 * time.Hour
	FileListTTL     = 30 * time.Minute
	FileExistsTTL   = 5 * time.Minute
	ProductListTTL  = 10 * time.Minute
	ProductTTL      = 30 * time.Minute
	UserProfileTTL  = 15 * time.Minute
)

func NewRedisCache(config config.CacheConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Redis connection failed", "error", err, "addr", fmt.Sprintf("%s:%d", config.Host, config.Port))
	} else {
		logger.Info("Redis connected successfully", "addr", fmt.Sprintf("%s:%d", config.Host, config.Port), "db", config.DB)
	}

	return &RedisCache{client: rdb}
}

func NewRedisCacheFromConfig(config CacheConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Redis connection failed", "error", err, "addr", config.Addr)
	} else {
		logger.Info("Redis connected successfully", "addr", config.Addr, "db", config.DB)
	}

	return &RedisCache{client: rdb}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		logger.ErrorCtx(ctx, "Redis SET marshal failed", "key", key, "error", err)
		return err
	}

	err = c.client.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis SET failed", "key", key, "error", err)
	} else {
		logger.DebugCtx(ctx, "Redis SET success", "key", key, "ttl", expiration)
	}
	return err
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.DebugCtx(ctx, "Redis cache miss", "key", key)
		} else {
			logger.ErrorCtx(ctx, "Redis GET failed", "key", key, "error", err)
		}
		return err
	}

	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		logger.ErrorCtx(ctx, "Redis GET unmarshal failed", "key", key, "error", err)
		return err
	}

	logger.DebugCtx(ctx, "Redis cache hit", "key", key)
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis DELETE failed", "key", key, "error", err)
	} else {
		logger.DebugCtx(ctx, "Redis DELETE success", "key", key)
	}
	return err
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis EXISTS failed", "key", key, "error", err)
		return false, err
	}

	exists := result > 0
	logger.DebugCtx(ctx, "Redis EXISTS check", "key", key, "exists", exists)
	return exists, nil
}

func (c *RedisCache) Ping(ctx context.Context) error {
	err := c.client.Ping(ctx).Err()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis PING failed", "error", err)
	}
	return err
}

func (c *RedisCache) Close() error {
	logger.Info("Closing Redis connection")
	return c.client.Close()
}

func (c *RedisCache) SetMultiple(ctx context.Context, data map[string]interface{}, expiration time.Duration) error {
	pipe := c.client.Pipeline()

	for key, value := range data {
		jsonData, err := json.Marshal(value)
		if err != nil {
			logger.ErrorCtx(ctx, "Redis SET MULTIPLE marshal failed", "key", key, "error", err)
			continue
		}
		pipe.Set(ctx, key, jsonData, expiration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Redis SET MULTIPLE failed", "error", err, "count", len(data))
	} else {
		logger.DebugCtx(ctx, "Redis SET MULTIPLE success", "count", len(data), "ttl", expiration)
	}
	return err
}

func (c *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	results, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis GET MULTIPLE failed", "error", err, "keys", keys)
		return nil, err
	}

	data := make(map[string]string)
	hitCount := 0

	for i, result := range results {
		if result != nil {
			if str, ok := result.(string); ok {
				data[keys[i]] = str
				hitCount++
			}
		}
	}

	logger.DebugCtx(ctx, "Redis GET MULTIPLE completed",
		"total", len(keys), "hits", hitCount, "misses", len(keys)-hitCount)

	return data, nil
}

func (c *RedisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetchFn func() (interface{}, error)) error {

	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	if err != redis.Nil {
		logger.ErrorCtx(ctx, "Redis GET OR SET - cache error", "key", key, "error", err)
	}

	data, err := fetchFn()
	if err != nil {
		logger.ErrorCtx(ctx, "Redis GET OR SET - fetch failed", "key", key, "error", err)
		return err
	}

	if setErr := c.Set(ctx, key, data, ttl); setErr != nil {
		logger.WarnCtx(ctx, "Redis GET OR SET - cache set failed", "key", key, "error", setErr)
	}

	jsonData, _ := json.Marshal(data)
	return json.Unmarshal(jsonData, dest)
}
