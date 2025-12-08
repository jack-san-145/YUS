package redis

import (
	"context"
	"fmt"
	"time"
)

func (r *RedisStore) GetOtp(ctx context.Context, email string) (string, error) {
	otp, err := r.RedisClient.Get(ctx, email).Result()
	if err != nil {
		fmt.Println("error while get the otp from redis - ", err)
		return "", err
	}
	return otp, nil
}

func (r *RedisStore) SetOtp(ctx context.Context, email string, otp string) error {
	err := r.RedisClient.Set(ctx, email, otp, 3*time.Minute).Err()
	if err != nil {
		fmt.Println("error while set the otp to redis - ", err)
		return err
	}
	return nil
}
