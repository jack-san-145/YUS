package driver

import (
	"yus/internal/storage"
)

type DriverHandler struct {
	Store *storage.Store
}

func NewDriverHandler(store *storage.Store) *DriverHandler {
	return &DriverHandler{
		Store: store,
	}
}
