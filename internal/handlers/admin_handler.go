package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/services"
	"yus/internal/storage/redis"
)

func Admin_otp_handler(w http.ResponseWriter, r *http.Request) {
	var admin_register_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	fmt.Println(email, name, password)

	if services.ValidateClgMail(email) && services.ValidateName(name) && services.ValidatePassword(password) {
		otp := services.GenerateOtp()
		go services.SendEmailTo(email, otp)
		admin_register_status["status"] = "Otp sent"
		// admin_register_status["status"] = redis.StoreAdmin(name, email, password)
	} else {
		admin_register_status["status"] = "invalid"
	}
	WriteJSON(w, r, admin_register_status)
}

func Verify_admin_otp(w http.ResponseWriter, r *http.Request) {
	var admin_register_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	fmt.Println("verify otp for - ", email, name, password)

	if services.ValidateClgMail(email) && services.ValidateName(name) && services.ValidatePassword(password) {
		// admin_register_status["status"] = redis.StoreAdmin(name, email, password)
	} else {
		admin_register_status["status"] = "invalid"
	}
	WriteJSON(w, r, admin_register_status)
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
