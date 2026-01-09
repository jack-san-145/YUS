package admin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/models"
	"yus/internal/services"
)

func (h *AdminHandler) GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	current_schedule, _ := h.Store.DB.GetCurrentSchedule(ctx)
	response.WriteJSON(w, r, current_schedule)

}

// To add many drivers in one request
func (h *AdminHandler) AddDriverHandler(w http.ResponseWriter, r *http.Request) {

	var (
		Driver_array []models.Driver
		Status_array []models.DriverAddedStatus
	)

	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&Driver_array)
	if err != nil {
		log.Println("error while decode the driver array - ", err)
		return
	}
	for _, driver := range Driver_array {
		var status models.DriverAddedStatus
		status.Name = driver.Name
		status.MobileNo = driver.Mobile_no

		if services.ValidateMobileNo(driver.Mobile_no) && services.ValidateName(driver.Name) {
			if err := h.Store.DB.AddDriver(ctx, &driver); err == nil { //stores the new_driver to DB
				status.IsAdded = true
			} else {
				status.IsAdded = false
			}

		} else {
			status.IsAdded = false
		}
		Status_array = append(Status_array, status)
	}
	response.WriteJSON(w, r, Status_array)

}

func (h *AdminHandler) ListDriversHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	//to load all the available routes
	all_available_drivers, _ := h.Store.DB.GetAvailableDrivers(ctx)
	response.WriteJSON(w, r, all_available_drivers)
}

func (h *AdminHandler) ScheduleBusHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	//https://yus.kwscloud.in/schedule-bus?bus_id=10&driver_id=1050&route_id=33

	var (
		schedule models.CurrentSchedule
		err      error
	)

	bus_id_string := r.URL.Query().Get("bus_id")
	driver_id_string := r.URL.Query().Get("driver_id")
	route_id_string := r.URL.Query().Get("route_id")

	schedule.BusId, err = strconv.Atoi(bus_id_string)
	if err != nil {
		log.Println("error while converting bus_id to int - ", err)
		return
	}

	schedule.DriverId, err = strconv.Atoi(driver_id_string)
	if err != nil {
		log.Println("error while converting driver_id to int - ", err)
		return
	}
	schedule.RouteId, err = strconv.Atoi(route_id_string)
	if err != nil {
		log.Println("error while converting route_id to int - ", err)
		return
	}

	err = h.Store.DB.ScheduleBus(ctx, &schedule)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"status": true})

	go func() {
		current_route, _ := h.Store.DB.GetCurrentBusRoutes(context.Background())
		h.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route)
	}()
}
