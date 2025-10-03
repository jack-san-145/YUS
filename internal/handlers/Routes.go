package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/storage/postgres"
)

func Load_all_available_routes(w http.ResponseWriter, r *http.Request) {

	//to load all the available routes
	all_available_routes := postgres.Load_available_routes()
	fmt.Println("avalaible routes - ", all_available_routes)
	WriteJSON(w, r, all_available_routes)
}

func Load_all_available_drivers(w http.ResponseWriter, r *http.Request) {

	//to load all the available routes
	all_available_drivers := postgres.Available_drivers()
	fmt.Println("avalaible drivers - ", all_available_drivers)
	WriteJSON(w, r, all_available_drivers)
}

func Add_New_Bus_handler(w http.ResponseWriter, r *http.Request) {

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

	//mapping driver to the bus
	var status = make(map[string]bool)
	driver_id_string := r.URL.Query().Get("driver_id")
	bus_id_string := r.URL.Query().Get("bus_id")

	driver_id_int, _ := strconv.Atoi(driver_id_string)
	bus_id_int, _ := strconv.Atoi(bus_id_string)
	fmt.Println(driver_id_int, bus_id_int)

	err := postgres.Map_driver_with_bus(driver_id_int, bus_id_int)
	if err != nil {
		status["driver_mapped"] = false
	} else {
		status["driver_mapped"] = true
	}
	WriteJSON(w, r, status)
}
