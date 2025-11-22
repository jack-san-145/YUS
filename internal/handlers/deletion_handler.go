package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

// to delele the given route from DB
func Remove_route_handler(w http.ResponseWriter, r *http.Request) {
	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	route_id_string := chi.URLParam(r, "route-id")
	route_id_int, err := strconv.Atoi(route_id_string)
	if err != nil {
		fmt.Println("error while converting route-id from string to int - ", err)
		WriteJSON(w, r, map[string]bool{"deleted": false})
		return
	}

	status := postgres.Remove_route(route_id_int)
	WriteJSON(w, r, status)

}

// to remove the given bus from DB
func Remove_bus_handler(w http.ResponseWriter, r *http.Request) {
	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	bus_id_string := chi.URLParam(r, "bus-id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting bus-id from string to int - ", err)
		WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}

	status := postgres.Remove_bus(bus_id_int)
	WriteJSON(w, r, status)

}

// to remove the given driver from DB
func Remove_driver_handler(w http.ResponseWriter, r *http.Request) {
	if !FindAdminSession_web(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	driver_id_string := chi.URLParam(r, "driver-id")
	driver_id_int, err := strconv.Atoi(driver_id_string)
	if err != nil {
		fmt.Println("error while converting driver-id from string to int - ", err)
		WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}

	status := postgres.Remove_driver(driver_id_int)
	WriteJSON(w, r, status)

}
