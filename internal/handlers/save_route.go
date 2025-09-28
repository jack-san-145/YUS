package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
	"yus/internal/storage/postgres"
)

var All_Bus_Routes []models.Route

func Save_New_route_handler(w http.ResponseWriter, r *http.Request) {
	var NewRoute models.Route
	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		fmt.Println("error while get the new route - ", err)
		return
	}
	NewRoute.Direction = "UP"
	postgres.SaveRoute_to_DB(&NewRoute)
	// All_Bus_Routes = append(All_Bus_Routes, NewRoute)
	// display_all_routes()

	/*
		// 	admin-app needs to send json like the below format:

			{
				"up_route_name": "Sattur to Madurai",
				"down_route_name":"Madurai to Sattur",
				"src":"Sattur",
				"dest":"Madurai",
				"stops": [
							{"location_name": "Sattur","lat": "9.3540035",  "lon": "77.9231079","is_stop":true,"arrival_time":"08:30"},
							{"location_name": "Virudhunagar", "lat": "9.3538361", "lon": "77.9231022","is_stop":true,"arrival_time":"09:15"},
							{"location_name": "Thirumangalam","lat": "9.3538282",  "lon": "77.9231091","is_stop":false,"arrival_time":"09:45"},
							{"location_name": "Madurai","lat": "9.353818",  "lon": "77.9231322","is_stop":true,"arrival_time":"10:50"}
						],
				"departure_time":"16:20",
			}
	*/
}

// func display_all_routes() {
// 	for _, bus_routes := range All_Bus_Routes {
// 		fmt.Println("bus route - ", bus_routes)
// 		fmt.Println("\n")
// 	}
// }
