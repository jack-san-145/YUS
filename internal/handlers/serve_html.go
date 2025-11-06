package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Serve_index_page(w http.ResponseWriter, r *http.Request) {

	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	templ, err := template.ParseFiles("../../ui/templates/index.html")
	if err != nil {
		fmt.Println("index.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_logo_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/templates/logo.html")
	if err != nil {
		fmt.Println("logo.html not found - ", err)
		return
	}

	// err=templ.Execute(w, nil)
	templ.Execute(w, nil)

}

func Serve_login_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/templates/login.html")
	if err != nil {
		fmt.Println("login.html not found - ", err)
		return
	}

	// err=templ.Execute(w, nil)
	templ.Execute(w, nil)

}

func Serve_otp_verify_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/templates/otpverify.html")
	if err != nil {
		fmt.Println("otpverify.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_register_page(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../../ui/templates/registerform.html")
	if err != nil {
		fmt.Println("registerform.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_bus_schedule_page(w http.ResponseWriter, r *http.Request) {

	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	templ, err := template.ParseFiles("../../ui/templates/bus_schedule.html")
	if err != nil {
		fmt.Println("bus_schedule.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}

func Serve_driver_page(w http.ResponseWriter, r *http.Request) {

	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	templ, err := template.ParseFiles("../../ui/templates/driver.html")
	if err != nil {
		fmt.Println("driver.html not found - ", err)
		return
	}
	templ.Execute(w, nil)

}
