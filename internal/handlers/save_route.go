package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
)

func Save_New_route_handler(w http.ResponseWriter, r *http.Request) {
	var new_route []models.RouteStops
	err := json.NewDecoder(r.Body).Decode(&new_route)
	if err != nil {
		fmt.Println("error while get the new route - ", err)
		return
	}
	fmt.Println("new route - ", new_route)
}
