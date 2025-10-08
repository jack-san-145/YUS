package redis

import (
	"context"
	"fmt"
	"yus/internal/services"
)

func StoreAdmin(name string, email string, password string) string {
	if Check_admin_exist() {
		fmt.Println("Admin already exists")
		return "Admin already exists"
	}
	hased_pass := services.Hash_this_password(password)
	err := rc.HSet(context.Background(), "Admin-data", "name", name, "email", email, "password", hased_pass).Err()
	if err != nil {
		fmt.Println("error while set the admin details to the redis - ", err)
		return "invalid"
	}
	return "successfully added admin"
}

func Check_admin_exist() bool {
	Exists, err := rc.Exists(context.Background(), "Admin-data").Result()
	if err != nil {
		fmt.Println("error while checking the existance of Admin-data - ", err)
		return false
	}
	if Exists == 1 {
		return true
	}
	return false

}

func Validate_Admin_login(email string, password string) (bool, string) {
	value, err := rc.HMGet(context.Background(), "Admin-data", "email", "password").Result() //to get the multiple values in a single query
	if err != nil {
		fmt.Println("error while accessing the Admin-data - ", err)
		return false, ""
	}

	if value[0] == nil || value[1] == nil {
		fmt.Println("both are nil")
		return false, ""
	}

	rc_email, ok1 := value[0].(string)
	rc_password, ok2 := value[1].(string)

	if !ok1 || !ok2 {
		return false, ""
	}

	if rc_email == email && services.Is_password_matched(rc_password, password) {
		return true, rc_email
	}
	return false, ""

}

func Change_Admin_password() {

}
