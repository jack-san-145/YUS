package main

import (
	"fmt"
	"net/http"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Get("/ws", handlers.Ws_hanler)
	fmt.Println("Server listening on :8090")
	err := http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}

}
