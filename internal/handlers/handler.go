package handlers

import (
	"net/http"
	"yus/internal/AppPkg"
	"yus/internal/storage"
)

type AdminHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	SendOTPHandler(w http.ResponseWriter, r *http.Request)
	VerifyOTPHandler(w http.ResponseWriter, r *http.Request)

	AddDriverHandler(w http.ResponseWriter, r *http.Request)
	RemoveDriverHandler(w http.ResponseWriter, r *http.Request)
	ListDriversHandler(w http.ResponseWriter, r *http.Request)

	AddBusHandler(w http.ResponseWriter, r *http.Request)
	RemoveBusHandler(w http.ResponseWriter, r *http.Request)
	AssignDriverToBusHandler(w http.ResponseWriter, r *http.Request)

	SaveRouteHandler(w http.ResponseWriter, r *http.Request)
	RemoveRouteHandler(w http.ResponseWriter, r *http.Request)
	ListRoutesHandler(w http.ResponseWriter, r *http.Request)
	GetCachedRoutesHandler(w http.ResponseWriter, r *http.Request)
	UpdateRouteDirectionHandler(w http.ResponseWriter, r *http.Request)
	AssignRouteToBusHandler(w http.ResponseWriter, r *http.Request)

	GetScheduleHandler(w http.ResponseWriter, r *http.Request)
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
