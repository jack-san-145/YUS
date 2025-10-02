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

func Map_Route_With_Bus_handler(w http.ResponseWriter, r *http.Request) {
	route_id_string := chi.URLParam(r, "route_id")
	bus_id_string := chi.URLParam(r, "bus_id")

	route_id_int, _ := strconv.Atoi(route_id_string)
	bus_id_int, _ := strconv.Atoi(bus_id_string)
	fmt.Println(route_id_int, bus_id_int)
}
