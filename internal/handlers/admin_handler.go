package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/services"
)

func Add_new_admin_handler(w http.ResponseWriter, r *http.Request) {
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
