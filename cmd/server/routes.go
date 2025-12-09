package main

import (
	"net/http"
	"yus/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(h *handlers.YUSHandler) *chi.Mux {

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

	/*



		admin-operations(website)

	*/
	router.Post("/yus/admin-login", h.Admin.LoginHandler)
	router.Post("/yus/send-otp-admin", h.Admin.SendOTPHandler)
	router.Post("/yus/verify-otp-admin", h.Admin.VerifyOTPHandler)
	router.Put("/yus/change-route/{direction}", h.Admin.UpdateRouteDirectionHandler)

	router.Post("/yus/add-new-driver", h.Admin.AddDriverHandler)
	router.Get("/yus/get-available-routes", h.Admin.ListRoutesHandler)

	//yus.kwscloud.in/yus/add-new-bus?bus_id=10
	router.Put("/yus/add-new-bus", h.Admin.AddBusHandler)

	//yus/allocate-bus-route?route_id=42&bus_id=10
	router.Put("/yus/allocate-bus-route", h.Admin.AssignRouteToBusHandler)
	router.Post("/yus/allocate-bus-driver", h.Admin.AssignDriverToBusHandler)
	router.Get("/yus/get-available-drivers", h.Admin.ListDriversHandler)

	//yus.kwscloud.in/yus/get-cached-routes?bus_id=10
	router.Get("/yus/get-cached-routes", h.Admin.GetCachedRoutesHandler)
	router.Get("/yus/get-current-schedule", h.Admin.GetScheduleHandler)

	//removal operations
	router.Delete("/yus/remove-route/{route-id}", h.Admin.RemoveRouteHandler)
	router.Delete("/yus/remove-bus/{bus-id}", h.Admin.RemoveBusHandler)
	router.Delete("/yus/remove-driver/{driver-id}", h.Admin.RemoveDriverHandler)

	/*



		admin-operations(mobile)

	*/
	router.Post("/yus/save-new-route", h.Admin.SaveRouteHandler)

	/*



		passenger-operations

	*/
	router.Get("/yus/passenger-ws", h.Passenger.WebSocketHandler)
	router.Get("/yus/src-{source}&dest-{destination}", h.Passenger.SrcDestHandler) //here i changed the endpoint format

	router.Get("/yus/src-{source}&dest-{destination}&stop-{stop}", h.Passenger.SrcDestStopsHandler)

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	router.Get("/yus/get-route", h.Passenger.GetRouteByBusIDHandler)

	router.Get("/yus/get-current-bus-routes", h.Passenger.GetCurrentBusRoutesHandler)

	/*



		driver-operations

	*/

	router.Post("/yus/send-otp-driver-password", h.Driver.SendOTPHandler)
	router.Post("/yus/verify-otp-driver-password", h.Driver.VerifyOTPHandler)
	router.Post("/yus/driver-login", h.Driver.LoginHandler)

	//wss://yus.kwscloud.in/yus/driver-ws?session_id='23sdf-sdfsq-341'
	router.Get("/yus/driver-ws", h.Driver.WebSocketHandler)

	//yus.kwscloud.in/yus/get-allotted-bus
	router.Get("/yus/get-allotted-bus", h.Driver.GetAllocatedBusHandler) //by sessions

	return router

}
