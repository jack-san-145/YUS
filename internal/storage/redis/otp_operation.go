package redis

import (
	"context"
	"fmt"
	"time"
)

func GetOtp(email string) string {
	otp, err := rc.Get(context.Background(), email).Result()
	if err != nil {
		fmt.Println("error while get the otp from redis - ", err)
		return ""
	}
	return otp
}

func SetOtp(email string, otp string) {
	err := rc.Set(context.Background(), email, otp, 3*time.Minute).Err()
	if err != nil {
		fmt.Println("error while set the otp to redis - ", err)
		return
	}
}
