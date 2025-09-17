package main

import (
	"fmt"
	"net/http"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Get("/driver-ws", handlers.Driver_Ws_hanler)
	router.Get("/passenger-ws,", handlers.Passenger_Ws_handler)
	fmt.Println("Server listening on :8090")
	err := http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}

}
