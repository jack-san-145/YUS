package passenger

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/services"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func (h *PassengerHandler) GetRouteByBusIDHandler(w http.ResponseWriter, r *http.Request) {

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting the bus_id_string to bus_id_int - ", err)
		response.WriteJSON(w, r, "null")
	}
	route_by_busId, _, _ := postgres.Find_route_by_bus_or_driver_ID(bus_id_int, "PASSENGER")
	fmt.Println("route_by bus id - ", route_by_busId)
	response.WriteJSON(w, r, route_by_busId)

}

func (h *PassengerHandler) SrcDestStopsHandler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	stop := chi.URLParam(r, "stop")

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	stop = services.Convert_to_CamelCase(stop)
	matched_routes := postgres.FindRoutes_by_src_dest_stop(src, dest, stop)
	response.WriteJSON(w, r, matched_routes)

}

func (h *PassengerHandler) SrcDestHandler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	fmt.Printf("given src - %v & destination - %v ", src, dest)

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	response.WriteJSON(w, r, postgres.FindRoutes_by_src_dest(src, dest))

}

func (h *PassengerHandler) GetCurrentBusRoutesHandler(w http.ResponseWriter, r *http.Request) {
	// bus_routes := postgres.Current_bus_routes()

	ctx := r.Context()
	bus_routes, err := h.Store.InMemoryDB.GetCachedRoute(ctx)

	if err != nil {
		bus_routes = postgres.Current_bus_routes()
		go h.Store.InMemoryDB.CacheBusRoute(ctx)
	}

	if len(bus_routes) != 0 {
		response.WriteJSON(w, r, bus_routes)
		return
	}
	response.WriteJSON(w, r, map[string]string{"bus_routes": "null"})
}
