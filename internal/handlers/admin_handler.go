package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/services"
	"yus/internal/storage/redis"
)

func Admin_registerhandler(w http.ResponseWriter, r *http.Request) {
	var admin_register_status map[string]string
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	if services.ValidateClgMail(email) && services.ValidateName(name) && services.ValidatePassword(password) {
		admin_register_status["status"] = "valid"
	} else {
		admin_register_status["status"] = "invalid"
	}
}

func Admin_login_handler(w http.ResponseWriter, r *http.Request) {
	var login_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing admin login form")
		login_status["login_status"] = "invalid"
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if redis.Validate_Admin_login(email, password) {
			login_status["login_status"] = "valid"
		} else {
			login_status["login_status"] = "invalid"
		}
	}
	WriteJSON(w, r, login_status)
}
