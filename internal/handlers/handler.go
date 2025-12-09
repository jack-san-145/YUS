package handlers

import (
	"yus/internal/AppPkg"
	"yus/internal/storage"
)

type AdminHandler interface {
}

type DriverHandler interface {
}

type PassengerHandler interface {
}

type YUSHandler struct {
	Store *storage.Store
	// it the store that contains both redis and postgres that can be used anywhere in the handlers package with "Store"

	Admin  AdminHandler
	Driver DriverHandler
	Pass   PassengerHandler
}

func NewHandler() *YUSHandler {
	return &YUSHandler{
		Store: AppPkg.NewApplication().Store,
	}
}
