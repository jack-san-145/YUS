package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"
)

// To add many drivers in one request
func Add_new_driver_handler(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession_web(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	var (
		Driver_array []models.Driver
		Status_array []models.DriverAddedStatus
	)
	err := json.NewDecoder(r.Body).Decode(&Driver_array)
	if err != nil {
		fmt.Println("error while decode the driver array - ", err)
		return
	}
	for _, driver := range Driver_array {
		var status models.DriverAddedStatus
		status.Name = driver.Name
		status.MobileNo = driver.Mobile_no

		fmt.Println("driver - ", driver)
		if services.ValidateMobileNo(driver.Mobile_no) && services.ValidateName(driver.Name) {
			if postgres.Store_new_driver_to_DB(&driver) { //stores the new_driver to DB
				status.IsAdded = true
			} else {
				status.IsAdded = false
			}

		} else {
			status.IsAdded = false
		}
		fmt.Println("status - ", status)
		Status_array = append(Status_array, status)
	}
	WriteJSON(w, r, Status_array)

}

func Load_all_available_drivers(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession_web(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//to load all the available routes
	all_available_drivers := postgres.Available_drivers()
	fmt.Println("avalaible drivers - ", all_available_drivers)
	WriteJSON(w, r, all_available_drivers)
}

func Driver_Otp_handler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var otp_status = make(map[string]bool)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing the form - ", err)
		return
	}
	driver_id := r.FormValue("driver_id")
	email := r.FormValue("email")

	driver_id_int, err := strconv.Atoi(driver_id)
	if err != nil {
		fmt.Println("error while converting the driver_id string to driver_id_int - ", err)
	}
	if !postgres.Check_Driver_exits(driver_id_int) {
		WriteJSON(w, r, map[string]string{"status": "no driver found"})
		return
	}
	if services.ValidateEmail(email) {
		otp := services.GenerateOtp()
		is_email_sent := services.SendEmailTo(email, otp)
		if is_email_sent {
			redis.SetOtp(ctx, email, otp) //set otp to redis if otp sent to email successfully
		}

		otp_status["otp_sent"] = is_email_sent
	} else {
		otp_status["otp_sent"] = false
	}

	WriteJSON(w, r, otp_status)
}

func Verify_driver_otp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var pass_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	driver_id_string := r.FormValue("driver_id")
	driver_id_int, err := strconv.Atoi(driver_id_string)
	if err != nil {
		fmt.Println("error while converting the driver_id string to driver_id_int - ", err)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	given_otp := r.FormValue("otp")
	fmt.Println("verify otp for - ", email, given_otp)

	if services.ValidateEmail(email) && services.ValidatePassword(password) {

		otp, _ := redis.GetOtp(ctx, email)
		if given_otp == otp {
			postgres.Set_driver_password(driver_id_int, email, password)
			pass_status["status"] = "success"

		} else {
			pass_status["status"] = "failed"
		}

	} else {
		pass_status["status"] = "failed"
	}
	WriteJSON(w, r, pass_status)
}

func Driver_login_handler(w http.ResponseWriter, r *http.Request) {
	var login_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing driver login form")
		login_status["login_status"] = "invalid"
	} else {
		driver_id := r.FormValue("driver_id")

		driver_id_int, err := strconv.Atoi(driver_id)
		if err != nil {
			fmt.Println("error while converting the driver_id string to driver_id_int - ", err)
		}
		password := r.FormValue("password")
		if postgres.ValidateDriver(driver_id_int, password) {
			login_status["login_status"] = "valid"
			session_id := redis.Create_Driver_Session(driver_id_int)
			login_status["session_id"] = session_id
		} else {
			login_status["login_status"] = "invalid"
		}
	}
	WriteJSON(w, r, login_status)
}

func Alloted_bus_handler(w http.ResponseWriter, r *http.Request) {

	isValid, driver_id := FindDriver_httpSession(r)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("driver_id- ", driver_id)
	alloted_bus := postgres.Get_Allotted_Bus(driver_id)
	if alloted_bus.BusID != 0 && alloted_bus.RouteId != 0 {
		WriteJSON(w, r, alloted_bus)
	} else {
		WriteJSON(w, r, map[string]string{"status": "no bus allotted"})
	}
}
