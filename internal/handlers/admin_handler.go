package handlers

import (
	"fmt"
	"net/http"
	"time"
	"yus/internal/services"
	"yus/internal/storage/redis"
)

func Admin_otp_handler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var admin_otp_status = make(map[string]any)
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

		exists, _ := redis.AdminExists(ctx)
		if exists {
			admin_otp_status["otp_sent"] = "Admin already exists"
			WriteJSON(w, r, admin_otp_status)
			return
		}

		//now synchronous
		otp := services.GenerateOtp()
		is_email_sent := services.SendEmailTo(email, otp)
		if is_email_sent {
			redis.SetOtp(email, otp) //set otp to redis if otp sent to email successfully
		}

		admin_otp_status["otp_sent"] = is_email_sent
	} else {
		admin_otp_status["otp_sent"] = false
	}
	WriteJSON(w, r, admin_otp_status)
}

func Verify_admin_otp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	//show this admin_register_status as it is on the frontend
	var admin_register_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")
	given_otp := r.FormValue("otp")
	fmt.Println("verify otp for - ", email, name, password, given_otp)

	if services.ValidateClgMail(email) && services.ValidateName(name) && services.ValidatePassword(password) {
		if given_otp == redis.GetOtp(email) {
			status, err := redis.AddAdmin(ctx, name, email, password)
			if err != nil {
				admin_register_status["status"] = status
			} else {
				admin_register_status["status"] = err.Error()
			}

		} else {
			admin_register_status["status"] = "invalid otp"
		}

	} else {
		admin_register_status["status"] = "invalid data"
	}
	WriteJSON(w, r, admin_register_status)
}

func Admin_login_handler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var login_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing admin login form")
		login_status["login_status"] = "invalid"
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")

		valid, _ := redis.AdminLogin(ctx, email, password)
		if valid {
			login_status["login_status"] = "valid"
			session_id := redis.Create_Admin_Session(email)
			login_status["session_id"] = session_id

			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    session_id,
				Path:     "/",
				Expires:  time.Now().Add(3 * time.Hour),
				HttpOnly: false,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, cookie)
		} else {
			login_status["login_status"] = "invalid"
		}
	}
	WriteJSON(w, r, login_status)
}
