package admin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/models"

	"github.com/go-chi/chi/v5"
)

func (h *AdminHandler) GetCachedRoutesHandler(w http.ResponseWriter, r *http.Request) {
	//yus.kwscloud.in/yus/get-cached-routes?bus_id=10

	ctx := r.Context()

	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		log.Println("error while converting bus_id_int to bus_id_string - ", err)
		return
	}
	cached_bus_routes, _ := h.Store.DB.GetCachedRoutesByBusID(ctx, bus_id_int)
	response.WriteJSON(w, r, cached_bus_routes)
}

func (h *AdminHandler) ListRoutesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	//to load all the available routes
	all_available_routes, _ := h.Store.DB.GetAvailableRoutes(ctx)
	response.WriteJSON(w, r, all_available_routes)
}

func (h *AdminHandler) AddBusHandler(w http.ResponseWriter, r *http.Request) {
	//yus.kwscloud.in/yus/add-new-bus?bus_id=10

	ctx := r.Context()
	//to add a new bus
	var status = make(map[string]string)
	bus_id_string := r.URL.Query().Get("bus_id")

	bus_id_int, _ := strconv.Atoi(bus_id_string)
	err := h.Store.DB.AddBus(ctx, bus_id_int)
	if err != nil {
		status["status"] = err.Error()
	} else {
		status["status"] = "success"
	}
	response.WriteJSON(w, r, status)
}

func (h *AdminHandler) AssignRouteToBusHandler(w http.ResponseWriter, r *http.Request) {
	//yus/allocate-bus-route?route_id=42&bus_id=10

	ctx := r.Context()

	//mapping route to the bus
	var status = make(map[string]bool)
	route_id_string := r.URL.Query().Get("route_id")
	bus_id_string := r.URL.Query().Get("bus_id")

	route_id_int, _ := strconv.Atoi(route_id_string)
	bus_id_int, _ := strconv.Atoi(bus_id_string)

	err := h.Store.DB.AssignRouteToBus(ctx, route_id_int, bus_id_int)
	if err != nil {
		status["route_mapped"] = false
	} else {
		status["route_mapped"] = true
	}
	response.WriteJSON(w, r, status)

	go func() {
		current_route, _ := h.Store.DB.GetCurrentBusRoutes(context.Background())
		h.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route)
	}()
}

func (h *AdminHandler) AssignDriverToBusHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	//gets the driver and bus allocation array , after allocated it returns the results
	var DriverAllocation_array []models.DriverAllocation

	//mapping driver to the bus
	var status = make(map[int]bool)
	err := json.NewDecoder(r.Body).Decode(&DriverAllocation_array)
	if err != nil {
		log.Println("error while decode the bus_and_drivers - ", err)
		return
	}

	for _, allcoate_driver := range DriverAllocation_array {
		err := h.Store.DB.AssignDriverToBus(ctx, allcoate_driver.DriverId, allcoate_driver.BusId)
		if err != nil {
			status[allcoate_driver.BusId] = false
		} else {
			status[allcoate_driver.BusId] = true
		}
	}

	response.WriteJSON(w, r, status)

	go func() {
		current_route, _ := h.Store.DB.GetCurrentBusRoutes(context.Background())
		h.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route)
	}()
}

func (h *AdminHandler) UpdateRouteDirectionHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	direction := chi.URLParam(r, "direction")
	if direction == "UP" || direction == "DOWN" {
		if ok, _ := h.Store.DB.ChangeRouteDirection(ctx, direction); ok {
			response.WriteJSON(w, r, map[string]bool{"changed": true})
		} else {
			response.WriteJSON(w, r, map[string]bool{"changed": false})
		}
	} else {
		response.WriteJSON(w, r, map[string]bool{"changed": false})
	}

	go func() {
		current_route, _ := h.Store.DB.GetCurrentBusRoutes(context.Background())
		h.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route)
	}()

}
