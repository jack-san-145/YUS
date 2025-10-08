package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"
)

func Add_new_driver_handler(w http.ResponseWriter, r *http.Request) {

	if !FindAdminSession(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
		status.Email = driver.Email

		fmt.Println("driver - ", driver)
		if services.ValidateMobileNo(driver.Mobile_no) && services.ValidateName(driver.Name) && services.ValidateEmail(driver.Email) {
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

	if !FindAdminSession(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//to load all the available routes
	all_available_drivers := postgres.Available_drivers()
	fmt.Println("avalaible drivers - ", all_available_drivers)
	WriteJSON(w, r, all_available_drivers)
}

// frontend -> backend: (post)

// {
// 	"email":"23ucs145@kamarajengg.edu.in",
// 	"password":"jack@145"
// }

// backend -> frontend
// {
// 	"login_status":"valid" or "Invalid"
// }
