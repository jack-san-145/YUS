package redis

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (r *RedisStore) GenerateSessionID(ctx context.Context) (string, error) {
	return uuid.New().String(), nil
}

func (r *RedisStore) CreateDriverSession(ctx context.Context, driverID int) (string, error) {
	session_id, _ := r.GenerateSessionID(ctx)
	err := r.RedisClient.HSet(ctx, session_id, "category", "DRIVER", "driver_id", driverID).Err()
	if err != nil {
		log.Println("error while setting the new driver session to the redis - ", err)
		return "", err
	}

	//set TTL for session
	err = r.RedisClient.Expire(ctx, session_id, 7*24*time.Hour).Err()
	if err != nil {
		log.Println("error while setting ttl to session - ", err)
		return "", err
	}
	return session_id, nil
}

func (r *RedisStore) CreateAdminSession(ctx context.Context, adminEmail string) (string, error) {
	session_id, _ := r.GenerateSessionID(ctx)
	var expiry = 3 * time.Hour
	err := r.RedisClient.HSet(ctx, session_id, "category", "ADMIN", "admin-email", adminEmail).Err()
	if err != nil {
		log.Println("error while setting the new admin session to the redis - ", err)
		return "", err
	}

	//set TTL for session
	err = r.RedisClient.Expire(ctx, session_id, expiry).Err()
	if err != nil {
		log.Println("error while setting ttl to session - ", err)
		return "", err
	}
	return session_id, nil
}

func (r *RedisStore) CheckAdminSession(ctx context.Context, sessionID string) (bool, error) {
	admin_email, err := r.RedisClient.HGet(ctx, sessionID, "admin-email").Result()
	if err != nil {
		log.Println("error while getting the session from redis - ", err)
		return false, err
	}
	if admin_email != "" {
		return true, nil
	}
	return false, nil

}

func (r *RedisStore) CheckDriverSession(ctx context.Context, sessionID string) (bool, int, error) {
	driver_id, err := r.RedisClient.HGet(ctx, sessionID, "driver_id").Result()
	if err != nil {
		log.Println("error while getting the session from redis - ", err)
	}
	if driver_id != "" {
		driver_id_int, err := strconv.Atoi(driver_id)
		if err != nil {
			log.Println("error while converting the driver id - ", err)
			return false, 0, err
		}
		return true, driver_id_int, nil
	}
	return false, 0, nil

}

func (r *RedisStore) DeleteSession(ctx context.Context, sessionID string) error {
	err := r.RedisClient.Del(ctx, sessionID).Err()
	if err != nil {
		log.Println("error while deleting the sessions - ", err)
		return err
	}
	return nil
}
