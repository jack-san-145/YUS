package main

import (
	"fmt"
	"net/http"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Get("/yus/admin-index-page", handlers.Serve_admin_index)
	router.Post("/yus/save-new-route", handlers.Save_New_route_handler)
	router.Get("/driver-ws", handlers.Driver_Ws_hanler)
	router.Get("/yus/passenger-ws", handlers.Passenger_Ws_handler)
	router.Get("/yus/src-{source}&dest{destination}", handlers.Src_Dest_handler)
	fmt.Println("Server listening on :8090")
	err := http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
