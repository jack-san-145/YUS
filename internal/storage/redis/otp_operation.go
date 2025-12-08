package redis

import (
	"context"
	"fmt"
	"time"
)

func GetOtp(ctx context.Context, email string) (string, error) {
	otp, err := rc.Get(ctx, email).Result()
	if err != nil {
		fmt.Println("error while get the otp from redis - ", err)
		return "", err
	}
	return otp, nil
}

func SetOtp(ctx context.Context, email string, otp string) error {
	err := rc.Set(ctx, email, otp, 3*time.Minute).Err()
	if err != nil {
		fmt.Println("error while set the otp to redis - ", err)
		return err
	}
	return nil
}
