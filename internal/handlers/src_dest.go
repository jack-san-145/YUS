package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/services"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func Get_rotue_by_busID(w http.ResponseWriter, r *http.Request) {
	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting the bus_id_string to bus_id_int - ", err)
		WriteJSON(w, r, "null")
	}
	route_by_busId := postgres.Find_route_by_busID(bus_id_int,"PASSENGER")
	WriteJSON(w, r, route_by_busId)

}

func Src_Dest_Stop_handler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	stop := chi.URLParam(r, "stop")

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	stop = services.Convert_to_CamelCase(stop)
	matched_routes := postgres.FindRoutes_by_src_dest_stop(src, dest, stop)
	WriteJSON(w, r, matched_routes)

}

func Src_Dest_handler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	fmt.Printf("given src - %v & destination - %v ", src, dest)

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	WriteJSON(w, r, postgres.FindRoutes_by_src_dest(src, dest))

}
