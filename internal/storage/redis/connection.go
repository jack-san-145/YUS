package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// var rc *redis.Client

type RedisStore struct {
	RedisClient *redis.Client
}

func NewRedisClient() *RedisStore {
	return &RedisStore{}
}

func (r *RedisStore) CreateClient(ctx context.Context) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:            os.Getenv("REDIS_SERVER_IP") + ":" + os.Getenv("REDIS_SERVER_PORT"), //yus instance ip
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("redis connection successfull")

	return nil

}
