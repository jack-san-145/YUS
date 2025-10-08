package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/services"
	"yus/internal/storage/redis"
)

func Admin_otp_handler(w http.ResponseWriter, r *http.Request) {
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

		if redis.Check_admin_exist() {
			admin_otp_status["otp_sent"] = "Admin already exists"
			WriteJSON(w, r, admin_otp_status)
			return
		}

		// ch := make(chan bool) //channel to store the otp is sent to email or not
		// otp := services.GenerateOtp()

		// go services.SendEmailTo(ch, email, otp) //pass also the channel

		// //this line wait until that go routine puts value to the ch
		// is_email_sent := <-ch //receives the email sent status from the bool channel 'ch'

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
			admin_register_status["status"] = redis.StoreAdmin(name, email, password)
		} else {
			admin_register_status["status"] = "invalid otp"
		}

	} else {
		admin_register_status["status"] = "invalid data"
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
		isValid, admin_email := redis.Validate_Admin_login(email, password)
		if isValid {
			login_status["login_status"] = "valid"
			redis.Create_Admin_Session(admin_email)
		} else {
			login_status["login_status"] = "invalid"
		}
	}
	WriteJSON(w, r, login_status)
}
