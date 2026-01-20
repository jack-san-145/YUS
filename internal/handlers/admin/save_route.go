package admin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	// "strconv"
	"yus/internal/handlers/common/response"
	"yus/internal/models"
)

// var All_Bus_Routes []models.Route

func (h *AdminHandler) SaveSameRouteHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var NewRoute models.Route
	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		log.Println("error while get the new route - ", err)
		return
	}
	NewRoute.Direction = "UP"
	status, _ := h.Store.DB.SaveRoute(ctx, &NewRoute)
	response.WriteJSON(w, r, map[string]string{"status": status})

	/*
		// 	admin-app needs to send json like the below format:

			{
				"up_route_name": "Sattur to Kamaraj-College",
				"down_route_name":"Kamaraj-College to Sattur",
				"src":"Sattur",
				"dest":"Kamaraj-College",
				"stops": [
							{"location_name": "Sattur","lat": "9.3540035",  "lon": "77.9231079","is_stop":true,"arrival_time":"08:00"},
							{"location_name": "rr nagar", "lat": "9.3538361", "lon": "77.9231022","is_stop":true,"arrival_time":"08:15"},
							{"location_name": "soolakrai", "lat": "9.3538361", "lon": "77.9231022","is_stop":false,"arrival_time":"08:35"},
							{"location_name": "collectrate", "lat": "9.3538361", "lon": "77.9231022","is_stop":true,"arrival_time":"08:45"},
							{"location_name": "Kamaraj-College","lat": "9.3538282",  "lon": "77.9231091","is_stop":false,"arrival_time":"09:00"}
						],
				"down_departure_time":"16:40"
			}
	*/
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
