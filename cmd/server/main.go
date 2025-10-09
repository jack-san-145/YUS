package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"net/http"
	"yus/internal/handlers"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	router := chi.NewRouter()

	// ✅ Add CORS middleware
	router.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{"*"}, // or specific: []string{"http://localhost:3000"}
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Get("/yus/admin-index-page", handlers.Serve_admin_index)

	router.Post("/yus/save-new-route", handlers.Save_New_route_handler)
	router.Get("/driver-ws", handlers.Driver_Ws_hanler)
	router.Post("/yus/add-new-driver", handlers.Add_new_driver_handler)
	router.Post("/yus/send-otp-admin", handlers.Admin_otp_handler)
	router.Post("/yus/verify-otp-admin", handlers.Verify_admin_otp)
	router.Post("/yus/admin-login", handlers.Admin_login_handler)
	router.Post("/yus/send-otp-admin", handlers.Admin_otp_handler)
	router.Get("/yus/passenger-ws", handlers.Passenger_Ws_handler)
	router.Get("/yus/src-{source}&dest-{destination}", handlers.Src_Dest_handler) //here i changed the endpoint format

	router.Get("/yus/src-{source}&dest-{destination}&stop-{stop}", handlers.Src_Dest_Stop_handler)

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	router.Get("/yus/get-route", handlers.Get_rotue_by_busID)

	router.Get("/yus/get-available-routes", handlers.Load_all_available_routes)

	//yus.kwscloud.in/yus/add-new-bus?bus_id=10
	router.Put("/yus/add-new-bus", handlers.Add_New_Bus_handler)

	//yus/allocate-bus-route?route_id=42&bus_id=10
	router.Put("/yus/allocate-bus-route", handlers.Map_Route_With_Bus_handler)

	router.Post("/yus/allocate-bus-driver", handlers.Map_Driver_With_Bus_handler)

	router.Get("/yus/get-available-drivers", handlers.Load_all_available_drivers)

	//yus.kwscloud.in/yus/get-cached-routes?bus_id=10
	router.Get("/yus/get-cached-routes", handlers.Cached_route_handler)

	//driver

	router.Post("/yus/send-otp-driver-password", handlers.Driver_Otp_handler)
	router.Post("/yus/verify-otp-driver-password", handlers.Verify_driver_otp)
	router.Post("/yus/driver-login", handlers.Driver_login_handler)

	//yus.kwscloud.in/yus/get-allotted-bus?driver_id=10
	router.Get("/yus/get-allotted-bus")

	postgres.Connect()        //make a connection to postgres
	redis.CreateRedisClient() //made a redis client

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
