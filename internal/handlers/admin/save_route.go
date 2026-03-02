package admin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"yus/internal/handlers/common/response"
	"yus/internal/models"
)

func (h *AdminHandler) SaveSameRouteHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var NewRoute models.Route
	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		log.Println("error while get the new route - ", err)
		return
	}
	NewRoute.Direction = "UP"

	copied_route := NewRoute
	copied_route.Stops = make([]models.RouteStops, len(NewRoute.Stops))
	copy(copied_route.Stops, NewRoute.Stops)

	routeID, status, _ := h.Store.DB.SaveRoute(ctx, &NewRoute)
	if routeID > 0 {
		copied_route.Id = routeID
		go h.Store.DB.StoreToBackupRoute(context.Background(), "SAME", &copied_route)
	}
	response.WriteJSON(w, r, map[string]string{"status": status})
}

func (h *AdminHandler) SaveDifferentRouteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var NewRoute models.Route

	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		log.Println("error while get the new route - ", err)
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}

	log.Printf("direction - %v & id - %v", NewRoute.Direction, NewRoute.Id)

	copied_route := NewRoute
	copied_route.Stops = make([]models.RouteStops, len(NewRoute.Stops))
	copy(copied_route.Stops, NewRoute.Stops)

	routeID, err := h.Store.DB.SaveDifferentPathRoute(ctx, &NewRoute)
	if err != nil {
		log.Println("error while saving different path route - ", err)
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	response.WriteJSON(w, r, map[string]any{"status": true, "route_id": routeID})

	copied_route.Id = routeID
	go h.Store.DB.StoreToBackupRoute(context.Background(), "DIFFERENT", &copied_route)
}
