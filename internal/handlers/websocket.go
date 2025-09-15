package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// var ConnMap sync.Map
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

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

func Ws_hanler(w http.ResponseWriter, r *http.Request) {
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
	var current_location Location
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
		fmt.Printf("lattitude - %s & longitude - %s", current_location.Latitude, current_location.Longitude)
		fmt.Println("\n\n")
	}
}
