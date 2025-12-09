package admin

import (
	"yus/internal/storage"
)

type AdminHandler struct {
	Store *storage.Store
}

func NewAdminHandler(store *storage.Store) *AdminHandler {
	return &AdminHandler{
		Store: store,
	}
}
