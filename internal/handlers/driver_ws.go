package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
	// "yus/internal/service"

	"github.com/gorilla/websocket"
)

// var ConnMap sync.Map

// func ws_hanler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("working")
// 	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	}}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("error while upgrading the websocket - ", err)
// 		return
// 	}

// 	defer conn.Close()
// 	fmt.Println("driver connected successfully ")
// 	var current_location Location
// 	for {
// 		_, loc, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("error while reading the websocket message - ", err)
// 			return
// 		}
// 		err = json.Unmarshal(loc, &current_location)
// 		if err != nil {
// 			fmt.Println("error while unmarshaling the location - ", err)
// 		}
// 		fmt.Printf("lattitude - %s & longitude - %s", current_location.Latitude, current_location.Longitude)
// 		fmt.Println("\n\n")
// 	}

// }

func Driver_Ws_hanler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("working")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}
	listen_for_location(conn)

}

func listen_for_location(conn *websocket.Conn) {
	defer conn.Close()
	fmt.Println("driver connected successfully ")
	var current_location models.Location
	for {
		_, loc, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error while reading the websocket message - ", err)
			return
		}
		err = json.Unmarshal(loc, &current_location)
		if err != nil {
			fmt.Println("error while unmarshaling the location - ", err)
		}

		send_location_to_passenger(&current_location)
		// fmt.Printf("lattitude - %s & longitude - %s & Speed - %s ", current_location.Latitude, current_location.Longitude, current_location.Speed)
		fmt.Println("\n\n")
		// service.Reverse_Geocoding(current_location)
		fmt.Printf("lattitude - %v & longitude - %v & speed - %v", current_location.Latitude, current_location.Longitude, current_location.Speed)
	}
}
