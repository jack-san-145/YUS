package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"yus/internal/models"

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
	var requested_bus_route models.PassengerWsRequest

	err := conn.ReadJSON(&requested_bus_route)
	if err != nil {
		fmt.Println("error reading the passenger ws message - ", err)
	} else {

		Add_PassConn(requested_bus_route.DriverId, conn) //store the passenger ws to corresponding driver ws

	}

}

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
