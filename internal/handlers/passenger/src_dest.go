package passenger

import (
	"context"
	"net/http"
	"strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/services"

	"github.com/go-chi/chi/v5"
)

func (h *PassengerHandler) GetRouteByBusIDHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	bus_id_string := r.URL.Query().Get("bus_id")
	bus_id_int, err := strconv.Atoi(bus_id_string)
	if err != nil {
		response.WriteJSON(w, r, "null")
	}
	route, _ := h.Store.DB.FindRouteByBusOrDriverID(ctx, bus_id_int, "PASSENGER")
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
	matched_routes, _ := h.Store.DB.FindRoutesBySrcDstStop(ctx, src, dest, stop)
	response.WriteJSON(w, r, matched_routes)

}

func (h *PassengerHandler) SrcDestHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	route, _ := h.Store.DB.FindRoutesBySrcDst(ctx, src, dest)
	response.WriteJSON(w, r, route)

}

func (h *PassengerHandler) GetCurrentBusRoutesHandler(w http.ResponseWriter, r *http.Request) {
	// bus_routes := postgres.Current_bus_routes()

	ctx := r.Context()
	bus_routes, err := h.Store.InMemoryDB.GetCachedRoute(ctx)

	if err != nil {
		bus_routes, _ = h.Store.DB.GetCurrentBusRoutes(ctx)
		go func() {
			current_route, _ := h.Store.DB.GetCurrentBusRoutes(context.Background())
			h.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route)
		}()
	}

	if len(bus_routes) != 0 {
		response.WriteJSON(w, r, bus_routes)
		return
	}
	response.WriteJSON(w, r, map[string]string{"bus_routes": "null"})
}
