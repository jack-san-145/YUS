package redis

import (
	"context"
	"fmt"
	"yus/internal/services"
)

func StoreAdmin(name string, email string, password string) {
	hased_pass := services.Hash_this_password(password)
	err := rc.HSet(context.Background(), "driver-name", name, "driver-email", email, "password", hased_pass).Err()
	if err != nil {
		fmt.Println("error while set the admin details to the redis - ", err)
		return
	}
}
