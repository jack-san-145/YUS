package main

import (
	"net/http"
	"yus/internal/AppPkg"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(app *AppPkg.Application, h *handlers.YUSHandler) *chi.Mux {

	router := chi.NewRouter()

	// CORS middleware
	router.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Serve static files for admin website (CSS, JS, images)
	fileServer := http.FileServer(http.Dir("../../ui/Admin-website/static"))
	router.Handle("/admin-static/*", http.StripPrefix("/admin-static/", fileServer))

	//web pages
	router.Get("/", handlers.Serve_logo_page)
	router.Get("/yus/serve-index-page", handlers.Serve_index_page)
	router.Get("/yus/serve-login-page", handlers.Serve_login_page)
	router.Get("/yus/serve-bus-schedule-page", handlers.Serve_bus_schedule_page)
	router.Get("/yus/serve-driver-page", handlers.Serve_driver_page)
	router.Get("/yus/serve-register-page", handlers.Serve_register_page)
	router.Get("/yus/serve-otp-verify-page", handlers.Serve_otp_verify_page)

	//Passenger Operations
	router.Group(func(passenger chi.Router) {
		router.Get("/yus/passenger-ws", h.Passenger.WebSocketHandler)
		router.Get("/yus/get-current-bus-routes", h.Passenger.GetCurrentBusRoutesHandler)
		router.Get("/yus/src-{source}&dest-{destination}", h.Passenger.SrcDestHandler) //here i changed the endpoint format
		router.Get("/yus/src-{source}&dest-{destination}&stop-{stop}", h.Passenger.SrcDestStopsHandler)
		router.Get("/yus/get-route", h.Passenger.GetRouteByBusIDHandler) //route by BusID

	})

	//Driver Operations
	router.Group(func(driver chi.Router) {
		router.Post("/yus/send-otp-driver-password", h.Driver.SendOTPHandler)
		router.Post("/yus/verify-otp-driver-password", h.Driver.VerifyOTPHandler)
		router.Post("/yus/driver-login", h.Driver.LoginHandler)
	})

	router.Group(func(protectedDriver chi.Router) {
		protectedDriver.Use(app.IsDriverAuthorized)
		protectedDriver.Get("/yus/driver-ws", h.Driver.WebSocketHandler)
		protectedDriver.Get("/yus/get-allotted-bus", h.Driver.GetAllocatedBusHandler)
	})

	//Admin Operations
	router.Group(func(admin chi.Router) {
		admin.Post("/yus/admin-login", h.Admin.LoginHandler)
		admin.Post("/yus/send-otp-admin", h.Admin.SendOTPHandler)
		admin.Post("/yus/verify-otp-admin", h.Admin.VerifyOTPHandler)
	})

	router.Group(func(protectedAdmin chi.Router) {
		// protectedAdmin.Use(app.IsAdminAuthorized)

		//route creation
		router.Post("/yus/save-new-route", h.Admin.SaveRouteHandler)

		//scheduling Operations
		protectedAdmin.Put("/yus/change-route/{direction}", h.Admin.UpdateRouteDirectionHandler)
		protectedAdmin.Post("/yus/add-new-driver", h.Admin.AddDriverHandler)
		protectedAdmin.Get("/yus/get-available-routes", h.Admin.ListRoutesHandler)
		protectedAdmin.Put("/yus/add-new-bus", h.Admin.AddBusHandler)
		protectedAdmin.Put("/yus/allocate-bus-route", h.Admin.AssignRouteToBusHandler)
		protectedAdmin.Post("/yus/allocate-bus-driver", h.Admin.AssignDriverToBusHandler)
		protectedAdmin.Get("/yus/get-available-drivers", h.Admin.ListDriversHandler)
		protectedAdmin.Get("/yus/get-cached-routes", h.Admin.GetCachedRoutesHandler)
		protectedAdmin.Get("/yus/get-current-schedule", h.Admin.GetScheduleHandler)

		//removal operations
		protectedAdmin.Delete("/yus/remove-route/{route-id}", h.Admin.RemoveRouteHandler)
		protectedAdmin.Delete("/yus/remove-bus/{bus-id}", h.Admin.RemoveBusHandler)
		protectedAdmin.Delete("/yus/remove-driver/{driver-id}", h.Admin.RemoveDriverHandler)

	})

	return router

}
