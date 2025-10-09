package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
)

var DriverMap sync.Map //key: *websocket.Conn , value : int
// key = driver Websocket conncetion & value =  driver_id

// Store a driver connection
func Add_DriverConn(driverId int, conn *websocket.Conn) {

	DriverMap.Store(conn, driverId)
}

// Remove a driver connection
func Remove_DriverConn(conn *websocket.Conn) {
	DriverMap.Delete(conn)
}

// Get connections
// func Get_DriverConns(driverId int) []*websocket.Conn {
// 	value, ok := PassengerMap.Load(driverId)
// 	if !ok {
// 		return nil
// 	}
// 	return value.([]*websocket.Conn)
// }
