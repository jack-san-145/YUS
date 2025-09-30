package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func Add_new_admin_handler(w http.ResponseWriter, r *http.Request) {
	var admin_register_status map[string]string
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	fmt.Println("email - ", email)
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !isMatch {
		admin_register_status["status"] = "invalid"
		return
	}
	isGmail := strings.Split(email, "@")
	if isGmail[1] != "kamarajengg.edu.in" {
		admin_register_status["status"] = "invalid"
		return
	}
}
