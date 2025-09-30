package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var rc *redis.Client

func CreateRedisClient() {
	rc = redis.NewClient(&redis.Options{
		Addr:            os.Getenv("YUS_INSTANCE_IP") + ":" + os.Getenv("REDIS_SERVER_PORT"), //yus instance ip
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("redis client connection faliure - ", err)
		return
	}
	fmt.Println("redis connection successfull")
}

func GiveRedisConnection() *redis.Client {
	return rc
}
