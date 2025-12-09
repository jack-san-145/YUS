package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/handlers/common/response"
	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"
)

func Get_Schedule_handler(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession_web(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	current_schedule := postgres.Get_Current_schedule()
	response.WriteJSON(w, r, current_schedule)

}

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
	response.WriteJSON(w, r, Status_array)

}

func Load_all_available_drivers(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession_web(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//to load all the available routes
	all_available_drivers := postgres.Available_drivers()
	fmt.Println("avalaible drivers - ", all_available_drivers)
	response.WriteJSON(w, r, all_available_drivers)
}
