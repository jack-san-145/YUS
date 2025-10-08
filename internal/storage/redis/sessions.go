package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func generate_sessionID() string {
	return uuid.New().String()
}

func Create_Driver_Session(driver_id int) string {
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
	return session_id
}

func Create_Admin_Session(admin_email string) string {
	session_id := generate_sessionID()
	var expiry = 3 * time.Hour
	err := rc.HSet(context.Background(), session_id, "category", "ADMIN", "admin-email", admin_email).Err()
	if err != nil {
		fmt.Println("error while setting the new admin session to the redis - ", err)
	}

	//set TTL for session
	err = rc.Expire(context.Background(), session_id, expiry).Err()
	if err != nil {
		fmt.Println("error while setting ttl to session - ", err)
	}
	return session_id
}

func Check_Admin_session(session_id string) bool {
	admin_email, err := rc.HGet(context.Background(), session_id, "admin-email").Result()
	if err != nil {
		fmt.Println("error while getting the session from redis - ", err)
	}
	if admin_email != "" {
		return true
	}
	return false

}

func Check_Driver_session(session_id string) (bool, int) {
	driver_id, err := rc.HGet(context.Background(), session_id, "driver_id").Result()
	if err != nil {
		fmt.Println("error while getting the session from redis - ", err)
	}
	if driver_id != "" {
		driver_id_int, err := strconv.Atoi(driver_id)
		if err != nil {
			fmt.Println("error while converting the driver id - ", err)
			return false, 0
		}
		return true, driver_id_int
	}
	return false, 0

}

func Del_session(session_id string) {
	err := rc.Del(context.Background(), session_id).Err()
	if err != nil {
		fmt.Println("error while deleting the sessions - ", err)
	}
}
