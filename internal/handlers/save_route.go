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

}

// func display_all_routes() {
// 	for _, bus_routes := range All_Bus_Routes {
// 		fmt.Println("bus route - ", bus_routes)
// 		fmt.Println("\n")
// 	}
// }
