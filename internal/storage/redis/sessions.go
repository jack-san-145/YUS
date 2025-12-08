package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GenerateSessionID(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

func CreateDriverSession(ctx context.Context, driverID int) (string, error) {
	session_id, _ := GenerateSessionID(ctx)
	err := rc.HSet(context.Background(), session_id, "category", "DRIVER", "driver_id", driverID).Err()
	if err != nil {
		fmt.Println("error while setting the new driver session to the redis - ", err)
		return "", err
	}

	//set TTL for session
	err = rc.Expire(context.Background(), session_id, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Println("error while setting ttl to session - ", err)
		return "", err
	}
	return session_id, nil
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

func CheckDriverSession(ctx context.Context, sessionID string) (bool, int, error) {
	driver_id, err := rc.HGet(context.Background(), sessionID, "driver_id").Result()
	if err != nil {
		fmt.Println("error while getting the session from redis - ", err)
	}
	if driver_id != "" {
		driver_id_int, err := strconv.Atoi(driver_id)
		if err != nil {
			fmt.Println("error while converting the driver id - ", err)
			return false, 0, err
		}
		return true, driver_id_int, nil
	}
	return false, 0, nil

}

func Del_session(session_id string) {
	err := rc.Del(context.Background(), session_id).Err()
	if err != nil {
		fmt.Println("error while deleting the sessions - ", err)
	}
}
