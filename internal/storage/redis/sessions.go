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

func Create_Admin_Session(admin_email string) {
	session_id := generate_sessionID()
	var expiry time.Duration
	err := rc.HSet(context.Background(), session_id, "category", "ADMIN", "admin_id", admin_email).Err()
	if err != nil {
		fmt.Println("error while setting the new admin session to the redis - ", err)
	}

	//set TTL for session
	err = rc.Expire(context.Background(), session_id, expiry).Err()
	if err != nil {
		fmt.Println("error while setting ttl to session - ", err)
	}
}

func Get_session(session_id string) string {
	id, err := rc.HGet(context.Background(), session_id, "driver_id").Result()
	if err != nil {
		fmt.Println("error while getting the driver session from redis - ", err)
	}
	return id

}

func Del_session(session_id string) {
	err := rc.Del(context.Background(), session_id).Err()
	if err != nil {
		fmt.Println("error while deleting the sessions - ", err)
	}
}
