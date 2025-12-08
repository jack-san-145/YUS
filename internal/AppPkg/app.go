package AppPkg

import (
	"yus/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	Port   string
	Router *chi.Mux
	Store  *storage.Store
}

var App *Application
