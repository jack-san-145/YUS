package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/models"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func Cached_route_handler(w http.ResponseWriter, r *http.Request) {
	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting bus_id_int to bus_id_string - ", err)
		return
	}
	cached_bus_routes := postgres.Load_cached_route(bus_id_int)
	WriteJSON(w, r, cached_bus_routes)
}

func Load_all_available_routes(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//to load all the available routes
	all_available_routes := postgres.Load_available_routes()
	fmt.Println("avalaible routes - ", all_available_routes)
	WriteJSON(w, r, all_available_routes)
}

func Add_New_Bus_handler(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//to add a new bus
	var status = make(map[string]string)
	bus_id_string := r.URL.Query().Get("bus_id")

	bus_id_int, _ := strconv.Atoi(bus_id_string)
	err := postgres.Add_new_bus(bus_id_int)
	if err != nil {
		status["status"] = err.Error()
	} else {
		status["status"] = "success"
	}
	WriteJSON(w, r, status)
}

func Map_Route_With_Bus_handler(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//mapping route to the bus
	var status = make(map[string]bool)
	route_id_string := r.URL.Query().Get("route_id")
	bus_id_string := r.URL.Query().Get("bus_id")

	route_id_int, _ := strconv.Atoi(route_id_string)
	bus_id_int, _ := strconv.Atoi(bus_id_string)
	fmt.Println(route_id_int, bus_id_int)

	err := postgres.Map_route_with_bus(route_id_int, bus_id_int)
	if err != nil {
		status["route_mapped"] = false
	} else {
		status["route_mapped"] = true
	}
	WriteJSON(w, r, status)
}

func Map_Driver_With_Bus_handler(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	//gets the driver and bus allocation array , after allocated it returns the results
	var DriverAllocation_array []models.DriverAllocation

	//mapping driver to the bus
	var status = make(map[int]bool)
	err := json.NewDecoder(r.Body).Decode(&DriverAllocation_array)
	if err != nil {
		fmt.Println("error while decode the bus_and_drivers - ", err)
		return
	}

	for _, allcoate_driver := range DriverAllocation_array {
		err := postgres.Map_driver_with_bus(allcoate_driver.DriverId, allcoate_driver.BusId)
		if err != nil {
			status[allcoate_driver.BusId] = false
		} else {
			status[allcoate_driver.BusId] = true
		}
	}

	WriteJSON(w, r, status)
}

func ChangeRoute_direction_handler(w http.ResponseWriter, r *http.Request) {
	direction := chi.URLParam(r, "direction")
	if direction == "UP" || direction == "DOWN" {
		if postgres.Change_route_direction(direction) {
			WriteJSON(w, r, map[string]bool{"changed": true})
		} else {
			WriteJSON(w, r, map[string]bool{"changed": false})
		}
	} else {
		WriteJSON(w, r, map[string]bool{"changed": false})
	}
}
