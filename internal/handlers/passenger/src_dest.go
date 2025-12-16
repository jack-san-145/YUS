package passenger

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/services"
	"yus/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func (h *PassengerHandler) GetRouteByBusIDHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		fmt.Println("error while converting the bus_id_string to bus_id_int - ", err)
		response.WriteJSON(w, r, "null")
	}
	route, _ := postgres.FindRouteByBusOrDriverID(ctx, bus_id_int, "PASSENGER")
	fmt.Println("route_by bus id - ", route.Currentroute)
	response.WriteJSON(w, r, route.Currentroute)

}

func (h *PassengerHandler) SrcDestStopsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	stop := chi.URLParam(r, "stop")

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	stop = services.Convert_to_CamelCase(stop)
	matched_routes, _ := postgres.FindRoutesBySrcDstStop(ctx, src, dest, stop)
	response.WriteJSON(w, r, matched_routes)

}

func (h *PassengerHandler) SrcDestHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	fmt.Printf("given src - %v & destination - %v ", src, dest)

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	route, _ := postgres.FindRoutesBySrcDst(ctx, src, dest)
	response.WriteJSON(w, r, route)

}

func (h *PassengerHandler) GetCurrentBusRoutesHandler(w http.ResponseWriter, r *http.Request) {
	// bus_routes := postgres.Current_bus_routes()

	ctx := r.Context()
	bus_routes, err := h.Store.InMemoryDB.GetCachedRoute(ctx)

	if err != nil {
		bus_routes, _ = postgres.Current_bus_routes()
		go h.Store.InMemoryDB.CacheBusRoute(context.Background())
	}

	if len(bus_routes) != 0 {
		response.WriteJSON(w, r, bus_routes)
		return
	}
	response.WriteJSON(w, r, map[string]string{"bus_routes": "null"})
}
