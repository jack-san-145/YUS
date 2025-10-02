package handlers

import (
	"fmt"
	"net/http"
	"yus/internal/storage/postgres"
)

func Load_all_available_routes(w http.ResponseWriter, r *http.Request) {
	all_available_routes := postgres.Load_available_routes()
	fmt.Println("avalaible routes - ", all_available_routes)
	WriteJSON(w, r, all_available_routes)
}
