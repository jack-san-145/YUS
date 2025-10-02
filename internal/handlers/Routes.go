package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func Load_all_available_routes(w http.ResponseWriter, r *http.Request) {
	all_available_routes := postgres.Load_available_routes()
	fmt.Println("avalaible routes - ", all_available_routes)
	WriteJSON(w, r, all_available_routes)
}

func Add_New_Bus_handler(w http.ResponseWriter, r *http.Request) {
	var status = make(map[string]bool)
	bus_no_string := chi.URLParam(r, "bus_no")

	bus_no_int, _ := strconv.Atoi(bus_no_string)
	err := postgres.Add_new_bus(bus_no_int)
	if err != nil {
		status["bus_added"] = false
	} else {
		status["bus_added"] = true
	}
	WriteJSON(w, r, status)
}

func Map_Route_With_Bus_handler(w http.ResponseWriter, r *http.Request) {
	var status = make(map[string]bool)
	route_id_string := chi.URLParam(r, "route_id")
	bus_id_string := chi.URLParam(r, "bus_id")

	route_id_int, _ := strconv.Atoi(route_id_string)
	bus_id_int, _ := strconv.Atoi(bus_id_string)
	fmt.Println(route_id_int, bus_id_int)

	err := postgres.Map_route_with_bus(route_id_int, bus_id_int)
	if err != nil {
		status["mapped"] = false
	} else {
		status["mapped"] = true
	}
	WriteJSON(w, r, status)
}
