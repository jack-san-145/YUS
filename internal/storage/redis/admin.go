package redis

import (
	"context"
	"fmt"
	"yus/internal/services"
)

func (r *RedisStore) AddAdmin(ctx context.Context, name string, email string, password string) (string, error) {

	exists, _ := r.AdminExists(ctx)
	if exists {
		fmt.Println("Admin already exists")
		return "", fmt.Errorf("admin already exists")
	}

	hased_pass := services.Hash_this_password(password)
	err := r.RedisClient.HSet(ctx, "Admin-data", "name", name, "email", email, "password", hased_pass).Err()
	if err != nil {
		fmt.Println("error while set the admin details to the redis - ", err)
		return "", fmt.Errorf("invalid")
	}
	return "successfully added admin", nil
}

func (r *RedisStore) AdminExists(ctx context.Context) (bool, error) {
	Exists, err := r.RedisClient.Exists(ctx, "Admin-data").Result()
	if err != nil {
		fmt.Println("error while checking the existance of Admin-data - ", err)
		return false, err
	}
	if Exists == 1 {
		return true, nil
	}
	return false, nil

}

func (r *RedisStore) AdminLogin(ctx context.Context, email string, password string) (bool, error) {

	value, err := r.RedisClient.HMGet(ctx, "Admin-data", "email", "password").Result() //to get the multiple values in a single query
	if err != nil {
		fmt.Println("error while accessing the Admin-data - ", err)
		return false, err
	}

	if value[0] == nil || value[1] == nil {
		fmt.Println("both are nil")
		return false, nil
	}

	rc_email, ok1 := value[0].(string)
	rc_password, ok2 := value[1].(string)

	if !ok1 || !ok2 {
		return false, nil
	}

	if rc_email == email && services.Is_password_matched(rc_password, password) {
		return true, nil
	}
	return false, nil

}

func (r *RedisStore) RemoveAdminSession(ctx context.Context, sessionID string) error {
	err := r.RedisClient.Del(ctx, sessionID).Err()
	return err
}
