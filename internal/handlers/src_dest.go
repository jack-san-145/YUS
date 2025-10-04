package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/services"
	"yus/internal/storage/postgres"

	// "yus/internal/models"

	"github.com/go-chi/chi/v5"
)

func Src_Dest_handler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	fmt.Printf("given src - %v & destination - %v ", src, dest)

	src = services.Convert_to_CamelCase(src)
	dest = services.Convert_to_CamelCase(dest)
	WriteJSON(w, r, postgres.FindRoutes_by_src_dest(src, dest))

}

// func FindRoutes_by_src_dest(src string, dest string) {
// 	var route []models.Route
// 	// for _, bus_route := range All_Bus_Routes {
// 	// 	if bus_route.Stops[0].LocationName == src && bus_route.Stops[len(bus_route.Stops)-1].LocationName == dest {
// 	// 		fmt.Println("src and dest matched")
// 	// 		route = append(route, bus_route)
// 	// 	}
// 	// }
// 	// return route

// }
