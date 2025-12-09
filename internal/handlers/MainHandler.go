package handlers

import (
	"yus/internal/handlers/admin"
	"yus/internal/handlers/driver"

	"yus/internal/handlers/passenger"
	"yus/internal/storage"
)

type YUSHandler struct {
	Store *storage.Store
	// it the store that contains both redis and postgres that can be used anywhere in the handlers package with "Store"

	Admin     *admin.AdminHandler
	Driver    *driver.DriverHandler
	Passenger *passenger.PassengerHandler
}

func NewHandler(store *storage.Store) *YUSHandler {
	return &YUSHandler{
		Store:     store,
		Admin:     admin.NewAdminHandler(store),
		Driver:    driver.NewDriverHandler(store),
		Passenger: passenger.NewPassengerHandler(store),
	}
}
