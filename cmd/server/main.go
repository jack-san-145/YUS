package main

import (
	"fmt"
	"net/http"
	// "time"
	"yus/internal/handlers"
	"yus/internal/storage/postgres"
	"yus/internal/storage/postgres/service"
	"yus/internal/storage/redis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	router := chi.NewRouter()

	// CORS middleware
	router.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// router.Use(TimeoutMiddleware(5 * time.Second)) // 5-second timeout

	// Serve static files (CSS, JS, images)
	fileServer := http.FileServer(http.Dir("../../ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

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
	router.Post("/yus/admin-login", handlers.Admin_login_handler)
	router.Post("/yus/send-otp-admin", handlers.Admin_otp_handler)
	router.Post("/yus/verify-otp-admin", handlers.Verify_admin_otp)
	router.Put("/yus/change-route/{direction}", handlers.ChangeRoute_direction_handler)

	router.Post("/yus/add-new-driver", handlers.Add_new_driver_handler)
	router.Get("/yus/get-available-routes", handlers.Load_all_available_routes)

	//yus.kwscloud.in/yus/add-new-bus?bus_id=10
	router.Put("/yus/add-new-bus", handlers.Add_New_Bus_handler)

	//yus/allocate-bus-route?route_id=42&bus_id=10
	router.Put("/yus/allocate-bus-route", handlers.Map_Route_With_Bus_handler)
	router.Post("/yus/allocate-bus-driver", handlers.Map_Driver_With_Bus_handler)
	router.Get("/yus/get-available-drivers", handlers.Load_all_available_drivers)

	//yus.kwscloud.in/yus/get-cached-routes?bus_id=10
	router.Get("/yus/get-cached-routes", handlers.Cached_route_handler)
	router.Get("/yus/get-current-schedule", handlers.Get_Schedule_handler)

	/*



		admin-operations(mobile)

	*/
	router.Post("/yus/admin-login", handlers.Admin_login_handler)
	router.Post("/yus/save-new-route", handlers.Save_New_route_handler)

	/*



		passenger-operations

	*/
	router.Get("/yus/passenger-ws", handlers.Passenger_Ws_handler)
	router.Get("/yus/src-{source}&dest-{destination}", handlers.Src_Dest_handler) //here i changed the endpoint format

	router.Get("/yus/src-{source}&dest-{destination}&stop-{stop}", handlers.Src_Dest_Stop_handler)

	//yus.kwscloud.in/yus/get-route?bus_id={bus_id}
	router.Get("/yus/get-route", handlers.Get_rotue_by_busID)

	router.Get("/yus/get-current-bus-routes", handlers.Get_Current_bus_routes_handler)

	/*



		driver-operations

	*/

	router.Post("/yus/send-otp-driver-password", handlers.Driver_Otp_handler)
	router.Post("/yus/verify-otp-driver-password", handlers.Verify_driver_otp)
	router.Post("/yus/driver-login", handlers.Driver_login_handler)

	//wss://yus.kwscloud.in/yus/driver-ws?session_id='23sdf-sdfsq-341'
	router.Get("/yus/driver-ws", handlers.Driver_Ws_hanler)

	//yus.kwscloud.in/yus/get-allotted-bus
	router.Get("/yus/get-allotted-bus", handlers.Alloted_bus_handler) //by sessions

	/*





	 */

	postgres.Connect()        //make a connection to postgres
	redis.CreateRedisClient() //made a redis client

	go service.Automate_route_scheduling() //change the route direction on-runtime
	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}

// func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.TimeoutHandler(next, timeout, "request timed out")
// 	}
// }
