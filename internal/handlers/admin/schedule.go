package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"
)

func (h *AdminHandler) GetScheduleHandler(w http.ResponseWriter, r *http.Request) {

	current_schedule := postgres.Get_Current_schedule()
	response.WriteJSON(w, r, current_schedule)

}

// To add many drivers in one request
func (h *AdminHandler) AddDriverHandler(w http.ResponseWriter, r *http.Request) {

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

func (h *AdminHandler) ListDriversHandler(w http.ResponseWriter, r *http.Request) {

	//to load all the available routes
	all_available_drivers := postgres.Available_drivers()
	fmt.Println("avalaible drivers - ", all_available_drivers)
	response.WriteJSON(w, r, all_available_drivers)
}

func ScheduleBusHandler(w http.ResponseWriter, r *http.Request) {

	//https://yus.kwscloud.in/schedule-bus?bus_id=10&driver_id=1050&route_id=33

	var (
		bus_id    int
		driver_id int
		route_id  int
		err       error
	)

	bus_id_string := r.URL.Query().Get("bus_id")
	driver_id_string := r.URL.Query().Get("driver_id")
	route_id_string := r.URL.Query().Get("route_id")

	bus_id, err = strconv.Atoi(bus_id_string)
	if err != nil {
		log.Println("error while converting bus_id to int - ", err)
		return
	}

	driver_id, err = strconv.Atoi(driver_id_string)
	if err != nil {
		log.Println("error while converting driver_id to int - ", err)
		return
	}
	route_id, err = strconv.Atoi(route_id_string)
	if err != nil {
		log.Println("error while converting route_id to int - ", err)
		return
	}

}
