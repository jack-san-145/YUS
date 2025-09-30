package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rc *redis.Client

func CreateRedisClient() {
	rc = redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("redis client connection faliure - ", err)
		return
	}
	fmt.Println("redis connection successfull", err)
}

func GiveRedisConnection() *redis.Client {
	return rc
}
