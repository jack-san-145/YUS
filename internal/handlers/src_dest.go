package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Src_Dest_handler(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "source")
	dest := chi.URLParam(r, "destination")
	fmt.Printf("given src - %v & destination - %v ", src, dest)
}
