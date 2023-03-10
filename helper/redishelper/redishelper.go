package redishelper

import (
	"context"
	"fmt"
	"go-redis-example/model"
	"go-redis-example/pkg/redis"
	"time"
)

type RedisHelper struct {
	RedisCache redis.RedisCache
}

// NewRedisCache creates an instance
func NewRedisHelper(addrs string) *RedisHelper {
	return &RedisHelper{
		RedisCache: redis.NewRedisCache(addrs),
	}
}

func (r *RedisHelper) SetDataExample(ctx context.Context, key string, merchantKey model.DataExample,
	expiration time.Duration) (err error) {
	err = r.RedisCache.Set(ctx, key, merchantKey, expiration)
	if err != nil {
		return fmt.Errorf("failed to set data example to Redis with err: %w", err)
	}
	return nil
}

func (r *RedisHelper) GetDataExample(ctx context.Context, key string) (dataEx *model.DataExample, err error) {
	dataEx = &model.DataExample{}
	err = r.RedisCache.Get(ctx, key, dataEx)
	if err != nil {
		return dataEx, fmt.Errorf("failed to get data example from Redis with err: %w", err)
	}
	return dataEx, nil
}

func (r *RedisHelper) CheckDataExampleExisted(ctx context.Context, key string) (isExisted bool, err error) {
	isExisted, err = r.RedisCache.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check key exist in Redis with err: %w", err)
	}
	return isExisted, nil
}

func (r *RedisHelper) SetNXDataExample(ctx context.Context, key string, merchantKey model.DataExample,
	expiration time.Duration) (result bool, err error) {
	result, err = r.RedisCache.SetNX(ctx, key, merchantKey, expiration)
	if err != nil {
		return false, fmt.Errorf("failed to set nx data example in Redis with err: %w", err)
	}
	return result, nil
}

func (r *RedisHelper) DeleteDataExample(ctx context.Context, key string) (err error) {
	err = r.RedisCache.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete data example in Redis with err: %w", err)
	}
	return nil
}
