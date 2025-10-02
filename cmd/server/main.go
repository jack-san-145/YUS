package main

import (
	"fmt"
	"net/http"
	"yus/internal/handlers"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}
	router := chi.NewRouter()
	router.Get("/yus/admin-index-page", handlers.Serve_admin_index)
	router.Get("/yus/get-available-routes", handlers.Load_all_available_routes)
	router.Post("/yus/save-new-route", handlers.Save_New_route_handler)
	router.Get("/driver-ws", handlers.Driver_Ws_hanler)
	router.Post("/yus/add-new-driver", handlers.Add_new_driver_handler)
	router.Post("/yus/send-otp-admin", handlers.Admin_otp_handler)
	router.Post("/yus/verify-otp-admin", handlers.Verify_admin_otp)
	router.Post("/yus/admin-login", handlers.Admin_login_handler)
	router.Post("/yus/send-otp-admin", handlers.Admin_otp_handler)
	router.Get("/yus/passenger-ws", handlers.Passenger_Ws_handler)
	router.Get("/yus/src-{source}&dest-{destination}", handlers.Src_Dest_handler) //here i changed the endpoint format

	router.Put("/yus/set{route_id}-{bus_id}", handlers.Map_Route_With_Bus_handler)

	postgres.Connect()        //make a connection to postgres
	redis.CreateRedisClient() //made a redis client

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
