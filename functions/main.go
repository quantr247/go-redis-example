package main

import (
	"context"
	"fmt"
	"go-redis-example/model"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

const addrRedis = "[PUT_YOUR_ADDRESS_HERE]:[PUT_YOUR_PORT_HERE]"

func main() {
	lambda.Start(func(ctx context.Context) (string, error) {
		if err := run(ctx, nil); err != nil {
			return "ERROR", fmt.Errorf("ERROR: %+v", err)
		}
		return "OK", nil
	})
}

func run(ctx context.Context, args map[string]interface{}) error {
	s := NewService(addrRedis)
	cacheKeyFirst := "cache_key_first"
	dataExample := model.DataExample{
		ID:       1,
		Name:     "first",
		IsActive: true,
	}
	testSetDataExample(ctx, s, cacheKeyFirst, dataExample)
	testGetDataExample(ctx, s, cacheKeyFirst)
	testCheckExisted(ctx, s, cacheKeyFirst)

	cacheKeySecond := "cache_key_second"
	testCheckExisted(ctx, s, cacheKeySecond)

	testSetNXDataExample(ctx, s, cacheKeyFirst, dataExample)
	testSetNXDataExample(ctx, s, cacheKeySecond, dataExample)
	return nil
}

func testSetDataExample(ctx context.Context, s service, cacheKey string, dataEx model.DataExample) {
	if err := s.redisCache.SetDataExample(ctx, cacheKey, dataEx, time.Duration(time.Duration(60).Seconds())); err != nil {
		panic(err)
	}
}

func testGetDataExample(ctx context.Context, s service, cacheKey string) {
	dataExFromRedis, err := s.redisCache.GetDataExample(ctx, cacheKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("get data example from Redis: ", dataExFromRedis)
}

func testCheckExisted(ctx context.Context, s service, cacheKey string) {
	isExisted, err := s.redisCache.CheckDataExampleExisted(ctx, cacheKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s isExisted: %v\n", cacheKey, isExisted)
}

func testSetNXDataExample(ctx context.Context, s service, cacheKey string, dataEx model.DataExample) {
	resultSetNX, err := s.redisCache.SetNXDataExample(ctx, cacheKey, dataEx, time.Duration(time.Duration(60).Seconds()))
	if err != nil {
		panic(err)
	}
	fmt.Printf("setnx %s has result: %v\n", cacheKey, resultSetNX)
}
