package handlers

import (
	"fmt"
	"net/http"
	"time"
	"yus/internal/models"
	"yus/internal/storage/postgres"

	"github.com/gorilla/websocket"
)

func Passenger_Ws_handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("working passenger app websocket upgradation")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}
	fmt.Println("passenger connected successfully ")

	go listen_passenger_message(conn)
	// Passenger_All_WS_connections = append(Passenger_All_WS_connections, conn)

}

func listen_passenger_message(conn *websocket.Conn) {
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

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ping error:", err)
				return
			}
		}
	}()
	var (
		old_requested_bus_route models.PassengerWsRequest
		requested_bus_route     models.PassengerWsRequest
	)

	for {
		fmt.Println("listing passenger")
		err := conn.ReadJSON(&requested_bus_route)
		if err != nil {
			fmt.Println("error reading the passenger ws message - ", err)
			return
		}

		fmt.Println("requested_bus_route - ", requested_bus_route)
		if postgres.Check_route_exits_for_pass_Ws(requested_bus_route) {
			fmt.Printf("old driver - %v and new driver - %v\n", old_requested_bus_route.DriverId, requested_bus_route.DriverId)
			Remove_PassConn(old_requested_bus_route.DriverId, conn)
			Add_PassConn(requested_bus_route.DriverId, conn)
			old_requested_bus_route = requested_bus_route
		}
	}
}

// func listen_passenger_message(conn *websocket.Conn) {
// 	var (
// 		old_requested_bus_route models.PassengerWsRequest
// 		requested_bus_route     models.PassengerWsRequest
// 	)

// 	for {
// 		err := conn.ReadJSON(&requested_bus_route)
// 		if err != nil {
// 			fmt.Println("error reading the passenger ws message - ", err)
// 		} else {

// 			if postgres.Check_route_exits_for_pass_Ws(requested_bus_route) {
// 				fmt.Printf("old driver - %v and new driver - %v - ", old_requested_bus_route.DriverId, requested_bus_route.DriverId)
// 				Remove_PassConn(old_requested_bus_route.DriverId, conn)
// 				Add_PassConn(requested_bus_route.DriverId, conn) //store the passenger ws to corresponding driver ws
// 				old_requested_bus_route = requested_bus_route
// 			}

// 		}
// 	}

// }

// func send_location_to_passenger(current_location *models.Location) {

// 	// defer conn.Close()
// 	for _, conn := range Passenger_All_WS_connections {
// 		conn.WriteJSON(*current_location)
// 	}

// 	// conn.WriteJSON()
// 	// for {
// 	// 	_, loc, err := conn.ReadMessage()
// 	// 	if err != nil {
// 	// 		fmt.Println("error while reading the websocket message - ", err)
// 	// 		return
// 	// 	}
// 	// 	err = json.Unmarshal(loc, &current_location)
// 	// 	if err != nil {
// 	// 		fmt.Println("error while unmarshaling the location - ", err)
// 	// 	}
// 	// 	// fmt.Printf("lattitude - %s & longitude - %s & Speed - %s ", current_location.Latitude, current_location.Longitude, current_location.Speed)
// 	// 	fmt.Println("\n\n")
// 	// 	// service.Reverse_Geocoding(current_location)
// 	// 	fmt.Printf("lattitude - %v & longitude - %v & speed - ", current_location.Latitude, current_location.Longitude, current_location.Speed)
// 	// }

// }
