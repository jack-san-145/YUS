package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func Serve_index_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/index.html")
	if err != nil {
		log.Println("index.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_logo_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/logo.html")
	if err != nil {
		log.Println("logo.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_login_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/login.html")
	if err != nil {
		log.Println("login.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_otp_verify_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/otpverify.html")
	if err != nil {
		log.Println("otpverify.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_register_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/registerform.html")
	if err != nil {
		log.Println("registerform.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_bus_schedule_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/bus_schedule.html")
	if err != nil {
		log.Println("bus_schedule.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_driver_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/YUS-Admin/templates/driver.html")
	if err != nil {
		log.Println("driver.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func ServePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("../../ui/Privacy-Policy/privacy-policy.html")
	if err != nil {
		log.Println("privacy-policy.html not found - ", err)
		return
	}
	templ.Execute(w, nil)
}
