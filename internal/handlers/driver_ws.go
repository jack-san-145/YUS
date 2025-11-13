package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// "sync"
	"yus/internal/models"
	"yus/internal/storage/postgres"

	// "yus/internal/service"

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

	Ongoing_route := postgres.Find_route_by_busID(driver_id, "DRIVER").Stops
	fmt.Println("ongping route - ", Ongoing_route)

	go listen_for_location(driver_id, conn)

}

func listen_for_location(driver_id int, conn *websocket.Conn) {
	defer func() {
		Remove_Driver_from_passengerMap(driver_id)
		conn.Close()
	}()
	fmt.Println("driver connected successfully")

	//bus_status
	// var Arrival_status map[int]string

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
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ping error:", err)
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

		Send_location_to_passenger(driver_id, current_location)
		fmt.Printf("latitude - %v & longitude - %v & speed - %v\n",
			current_location.Latitude, current_location.Longitude, current_location.Speed)
	}
}
