// package handlers

// import (
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// var DriverMap sync.Map //key: *websocket.Conn , value : int
// // key = driver Websocket conncetion & value =  driver_id

// // Store a connection
// func Add_DriverConn(driverId int, conn *websocket.Conn) {

// 	DriverMap.Store(conn, driverId)
// }

// // Get connections
// func Get_DriverConns(driverId int) []*websocket.Conn {
// 	value, ok := PassengerMap.Load(driverId)
// 	if !ok {
// 		return nil
// 	}
// 	return value.([]*websocket.Conn)
// }

// // Remove a connection
// func Remove_DriverConn(driverId int, conn *websocket.Conn) {
// 	value, ok := PassengerMap.Load(driverId)
// 	if !ok {
// 		return
// 	}
// 	conns := value.([]*websocket.Conn)
// 	for i, ws_conn := range conns {
// 		if ws_conn == conn {
// 			conns = append(conns[:i], conns[i+1:]...)
// 			break
// 		}
// 	}
// 	PassengerMap.Store(driverId, conns)
// }
