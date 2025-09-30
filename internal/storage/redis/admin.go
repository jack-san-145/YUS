package redis

import (
	"context"
	"fmt"
	"yus/internal/services"
)

func StoreAdmin(name string, email string, password string) {
	hased_pass := services.Hash_this_password(password)
	err := rc.HSet(context.Background(), "Admin-data", "name", name, "email", email, "password", hased_pass).Err()
	if err != nil {
		fmt.Println("error while set the admin details to the redis - ", err)
		return
	}
}

func Validate_Admin_login(email string, password string) bool {
	password = services.Hash_this_password(password)
	value, err := rc.HMGet(context.Background(), "Admin-data", email, password).Result() //to get the multiple values in a single query
	if err != nil {
		fmt.Println("error while accessing the Admin-data - ", err)
		return false
	}
	rc_email := value[0].(string)
	rc_password := value[1].(string)

	if rc_email == email && services.Is_password_matched(rc_password, password) {
		return true
	}
	return false

}
