package main

import (
	"net/http"
	"time"
	"yus/internal/AppPkg"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router.Use(middleware.RequestID) //Assigns a unique ID to each request
	router.Use(middleware.RealIP)    // get correct client IP behind proxies
	// router.Use(middleware.Logger)                    // logs requests
	router.Use(middleware.Recoverer)                 // recovers from panics
	router.Use(middleware.Timeout(60 * time.Second)) // set request timeout
	router.Use(app.RateLimit)                        //limits requests from same ip

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
		passenger.Get("/yus/passenger-ws", h.Passenger.WebSocketHandler)
		passenger.Get("/yus/get-current-bus-routes", h.Passenger.GetCurrentBusRoutesHandler)
		passenger.Get("/yus/src-{source}&dest-{destination}", h.Passenger.SrcDestHandler) //here i changed the endpoint format
		passenger.Get("/yus/src-{source}&dest-{destination}&stop-{stop}", h.Passenger.SrcDestStopsHandler)
		passenger.Get("/yus/get-route", h.Passenger.GetRouteByBusIDHandler) //route by BusID

	})

	//Driver Operations
	router.Group(func(driver chi.Router) {
		driver.Post("/yus/send-otp-driver-password", h.Driver.SendOTPHandler)
		driver.Post("/yus/verify-otp-driver-password", h.Driver.VerifyOTPHandler)
		driver.Post("/yus/driver-login", h.Driver.LoginHandler)
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
		protectedAdmin.Post("/yus/save-new-route", h.Admin.SaveRouteHandler)

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
		protectedAdmin.Put("/yus/schedule-bus", h.Admin.ScheduleBusHandler)
		//removal operations
		protectedAdmin.Delete("/yus/remove-route/{route-id}", h.Admin.RemoveRouteHandler)
		protectedAdmin.Delete("/yus/remove-bus/{bus-id}", h.Admin.RemoveBusHandler)
		protectedAdmin.Delete("/yus/remove-driver/{driver-id}", h.Admin.RemoveDriverHandler)

	})

	return router

}
