package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func generate_sessionID() string {
	return uuid.New().String()
}

func Create_Driver_Session(driver_id int) {
	session_id := generate_sessionID()
	err := rc.HSet(context.Background(), session_id, "category", "DRIVER", "driver_id", driver_id).Err()
	if err != nil {
		fmt.Println("error while setting the new driver session to the redis - ", err)
	}

	//set TTL for session
	err = rc.Expire(context.Background(), session_id, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Println("error while setting ttl to session - ", err)
	}
}

func Get_Driver_session(session_id string) string {
	driver_id, err := rc.HGet(context.Background(), session_id, "driver_id").Result()
	if err != nil {
		fmt.Println("error while getting the driver session from redis - ", err)
	}
	return driver_id

}
