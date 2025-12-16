package admin

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/storage/postgres"
)

// to delele the given route from DB
func (h *AdminHandler) RemoveRouteHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	route_id_string := chi.URLParam(r, "route-id")
	route_id_int, err := strconv.Atoi(route_id_string)
	if err != nil {
		fmt.Println("error while converting route-id from string to int - ", err)
		response.WriteJSON(w, r, map[string]bool{"deleted": false})
		return
	}

	err = postgres.RemoveRoute(ctx, route_id_int)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"deleted": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"deleted": true})

}

// to remove the given bus from DB
func (h *AdminHandler) RemoveBusHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	bus_id_string := chi.URLParam(r, "bus-id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting bus-id from string to int - ", err)
		response.WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}

	err = postgres.RemoveBus(ctx, bus_id_int)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"removed": true})

}

// to remove the given driver from DB
func (h *AdminHandler) RemoveDriverHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	driver_id_string := chi.URLParam(r, "driver-id")
	driver_id_int, err := strconv.Atoi(driver_id_string)
	if err != nil {
		fmt.Println("error while converting driver-id from string to int - ", err)
		response.WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}

	err = postgres.RemoveDriver(ctx, driver_id_int)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"removed": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"removed": true})

}
