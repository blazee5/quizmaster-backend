package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		panic(err)
	}

	return client
}
