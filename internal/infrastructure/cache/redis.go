package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
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
