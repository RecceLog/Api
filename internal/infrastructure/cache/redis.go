package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		Protocol: 2,
	})

	err := rdb.Set(ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		return nil, fmt.Errorf("error setting value in redis: %s", err)
	}

	_, err = rdb.Get(ctx, "foo").Result()
	if err != nil {
		return nil, fmt.Errorf("error getting value in redis: %s", err)
	}

	return rdb, nil
}
