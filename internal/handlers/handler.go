package handlers

import (
	"yus/internal/AppPkg"
	"yus/internal/handlers/admin"
	"yus/internal/handlers/driver"

	"yus/internal/handlers/passenger"
	"yus/internal/storage"
)

type YUSHandler struct {
	Store *storage.Store
	// it the store that contains both redis and postgres that can be used anywhere in the handlers package with "Store"

	Admin  *admin.AdminHandler
	Driver *driver.DriverHandler
	Pass   *passenger.PassengerHandler
}

func NewHandler() *YUSHandler {
	return &YUSHandler{
		Store: AppPkg.NewApplication().Store,
	}
}
