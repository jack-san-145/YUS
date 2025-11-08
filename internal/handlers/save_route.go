package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
	"yus/internal/storage/postgres"
)

// var All_Bus_Routes []models.Route

func Save_New_route_handler(w http.ResponseWriter, r *http.Request) {

	if !FindAdminSession_mobile(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var NewRoute models.Route
	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		fmt.Println("error while get the new route - ", err)
		return
	}
	NewRoute.Direction = "UP"
	fmt.Printf("actual route - %+v", NewRoute)
	status := postgres.SaveRoute_to_DB(&NewRoute)
	WriteJSON(w, r, status)
	// All_Bus_Routes = append(All_Bus_Routes, NewRoute)
	// display_all_routes()

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

/* //function to provide the new routes

async function send_route_to_server(new_route) {
	const response = await fetch("https://yus.kwscloud.in/yus/save-new-route", {
	method: "POST",
	headers: { "Content-Type": "application/json" },
	body: JSON.stringify(new_route),
	});
}
*/

// func display_all_routes() {
// 	for _, bus_routes := range All_Bus_Routes {
// 		fmt.Println("bus route - ", bus_routes)
// 		fmt.Println("\n")
// 	}
// }
