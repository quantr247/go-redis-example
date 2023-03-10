package main

import (
	"context"
	"go-redis-example/helper/redishelper"
	"go-redis-example/model"
	"time"
)

type service struct {
	redisCache DataExampleRedisCache
}

// NewService create new instance.
func NewService(addrs string) service {
	redisHelper := redishelper.NewRedisHelper(addrs)
	return service{
		redisCache: redisHelper,
	}
}

// This interface in consumer side because package main just want to use 4 method.
type DataExampleRedisCache interface {
	SetDataExample(ctx context.Context, key string, dataEx model.DataExample, expiration time.Duration) (err error)
	GetDataExample(ctx context.Context, key string) (dataEx *model.DataExample, err error)
	CheckDataExampleExisted(ctx context.Context, key string) (isExisted bool, err error)
	SetNXDataExample(ctx context.Context, key string, dataEx model.DataExample, expiration time.Duration) (result bool, err error)
}
