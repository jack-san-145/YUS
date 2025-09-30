package redis

import (
	"context"
)

func StoreAdmin(name string, email string, password string) {

	err := rc.HSet(context.Background(), "driver-name", name, "driver-email", email, "password", password).Err()
}
