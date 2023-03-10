package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	go_redis "github.com/redis/go-redis/v9"
)

// RedisCache is the interface implemented for Redis.
// This interface in producer side because it will be expose for external service.
type RedisCache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, valueType interface{}) error
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}

type redisCache struct {
	client *go_redis.Client
}

// NewRedisCache creates an instance.
func NewRedisCache(addrs string) RedisCache {
	client, err := initRedis(addrs)
	if err != nil {
		panic(fmt.Errorf("failed to init redis with err: %w", err))
	}
	return &redisCache{
		client: client,
	}
}

func initRedis(addr string) (*go_redis.Client, error) {
	client := go_redis.NewClient(&go_redis.Options{
		Addr: addr,
		DB:   0,
	})
	_, err := client.Ping(context.Background()).Result()
	return client, err
}

// Exists to check key has been already existed.
//
// Return: true: is existed, false: not exist.
func (h *redisCache) Exists(ctx context.Context, key string) (isExisted bool, err error) {
	result, err := h.client.Exists(ctx, key).Result()
	if err != nil {
		return true, err
	}

	// result = 0 is means key does not exist
	if result == 0 {
		return false, nil
	}
	return true, nil
}

// Set to save key and value.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.client.Set(ctx, key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

// Get to get value by key. Need to pass pointer of valueType.
func (r *redisCache) Get(ctx context.Context, key string, valueType interface{}) (err error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &valueType)
	if err != nil {
		return err
	}

	return nil
}

// Del to delete key.
func (h *redisCache) Del(ctx context.Context, key string) (err error) {
	_, err = h.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Expire to set more expiration for the key.
func (h *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	_, err = h.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

// SetNX to set key and value.
// Return the result:
//
//	true: if key does not exist, SetNX is successfully.
//	false: if key has been already existed, SetNX failed
func (h *redisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (result bool, err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return h.client.SetNX(ctx, key, string(data), expiration).Result()
}
