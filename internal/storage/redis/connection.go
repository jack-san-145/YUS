package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// var rc *redis.Client

type RedisClient struct {
	RC *redis.Client
}

func NewRedisClient() *RedisClient {
	return &RedisClient{}
}

func (r *RedisClient) CreateClient(ctx context.Context) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:            os.Getenv("REDIS_SERVER_IP") + ":" + os.Getenv("REDIS_SERVER_PORT"), //yus instance ip
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		fmt.Println("redis connection faliure - ", err)
		return nil, err
	}
	fmt.Println("redis connection successfull")

	return redisClient, nil

}
