package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"

	"github.com/gorilla/websocket"
)

func Driver_Ws_hanler(w http.ResponseWriter, r *http.Request) {

	isValid, driver_id := FindDriver_wssSession(r)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("working")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}

	Add_Driver_to_passengerMap(driver_id)

	go listen_for_location(driver_id, conn)

}

func listen_for_location(driver_id int, conn *websocket.Conn) {
	//bus_status
	var Arrival_status = make(map[int]string)
	var done = make(chan struct{}) // using to shutdown the unwanted ticker goroutine
	// this struct{} doesn't create any buffer and assign value , so simply the select will wait to receive the value on done bcz now done is empty or zero , it doesn't consider struct{} as a value

	defer func() {
		close(done) // to close the ticker goroutine when the driver disconnects
		redis.Store_ArrivalStatus(driver_id, Arrival_status)
		Remove_Driver_from_passengerMap(driver_id)
		conn.Close()
	}()

	redis_as, err := redis.Get_ArrivalStatus(driver_id)
	if err == nil {
		Arrival_status = redis_as
	}

	Ongoing_route := postgres.Find_route_by_bus_or_driver_ID(driver_id, "DRIVER").Stops
	fmt.Println("ongoing route - ", Ongoing_route)

	fmt.Println("driver connected successfully")

	// Ping/pong settings
	const (
		pongWait   = 60 * time.Second
		pingPeriod = 50 * time.Second
		writeWait  = 10 * time.Second
	)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start ping ticker
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Println("ping error:", err)
					return
				}
			case <-done: // this runs immediately when the done channel gets closed
				return
			}
		}

	}()

	var current_location models.Location
	for {
		_, loc, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error while reading the driver's websocket message - ", err)
			return
		}
		err = json.Unmarshal(loc, &current_location)
		if err != nil {
			fmt.Println("error while unmarshaling the location - ", err)
			continue
		}

		stop_sequence, reached_time, is_reached := services.FindNearestStop(current_location.Latitude, current_location.Longitude, Ongoing_route)
		_, ok := Arrival_status[stop_sequence] //returns true only if the key exists otherwise returns false
		if is_reached && !ok {
			Arrival_status[stop_sequence] = reached_time
			fmt.Println("arrival status - ", Arrival_status)
			current_location.ArrivalStatus = Arrival_status
		}
		Send_location_to_passenger(driver_id, current_location)
		fmt.Printf("latitude - %v & longitude - %v & speed - %v\n",
			current_location.Latitude, current_location.Longitude, current_location.Speed)
	}
}
