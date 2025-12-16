package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yus/internal/handlers/passenger"
	"yus/internal/models"
	"yus/internal/services"
	"yus/internal/storage/postgres"

	"github.com/gorilla/websocket"
)

func (h *DriverHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	//wss://yus.kwscloud.in/yus/driver-ws?session_id='23sdf-sdfsq-341'

	driver_id := r.Context().Value("DRIVER_ID").(int)

	fmt.Println("working")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}

	//check if the driver exists on the passenger map
	ok := passenger.PassengerConnStore.DriverExists(driver_id)
	if !ok {

		//if driver doesn't exists add him to the passengerMap
		passenger.PassengerConnStore.AddDriver(driver_id)
	}

	go h.listenForLocation(driver_id, conn)

}

func (h *DriverHandler) listenForLocation(driver_id int, conn *websocket.Conn) {
	//bus_status
	var Arrival_status = make(map[int]string)
	var done = make(chan struct{}) // using to shutdown the unwanted ticker goroutine
	// this struct{} doesn't create any buffer and assign value , so simply the select will wait to receive the value on done bcz now done is empty or zero , it doesn't consider struct{} as a value

	defer func() {
		close(done) // to close the ticker goroutine when the driver disconnects
		h.Store.InMemoryDB.StoreArrivalStatus(context.Background(), driver_id, Arrival_status)
		passenger.PassengerConnStore.RemoveDriver(driver_id)
		conn.Close()
	}()

	redis_as, err := h.Store.InMemoryDB.GetArrivalStatus(context.Background(), driver_id)
	if err == nil {
		Arrival_status = redis_as
	}

	r, _ := postgres.FindRouteByBusOrDriverID(context.Background(), driver_id, "DRIVER")
	Ongoing_route_stops := r.Currentroute.Stops
	fmt.Println("ongoing route - ", Ongoing_route_stops)

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

		stop_sequence, reached_time, is_reached := services.FindNearestStop(current_location.Latitude, current_location.Longitude, Ongoing_route_stops)
		_, ok := Arrival_status[stop_sequence] //returns true only if the key exists otherwise returns false
		if is_reached && !ok {
			Arrival_status[stop_sequence] = reached_time
			fmt.Println("arrival status - ", Arrival_status)
			current_location.ArrivalStatus = Arrival_status
		}
		passenger.PassengerConnStore.BroadcastLocation(driver_id, current_location)
		fmt.Printf("latitude - %v & longitude - %v & speed - %v\n",
			current_location.Latitude, current_location.Longitude, current_location.Speed)
	}
}
