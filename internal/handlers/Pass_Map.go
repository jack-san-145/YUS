package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
)



var PassengerMap sync.Map // key: int, value: []*websocket.Conn

// key = driver_id & value = passenger Websocket conncetions

// Store a connection
func Add_PassConn(driverId int, conn *websocket.Conn) {
	value, _ := PassengerMap.LoadOrStore(driverId, []*websocket.Conn{})
	conns := value.([]*websocket.Conn)  //assings the passenger ws connections to the conns
	conns = append(conns, conn)         //add the new passenger to the existing array and it mapped with driverId
	PassengerMap.Store(driverId, conns) //store the final conns to the ConnMap
}

// Get connections
func Get_PassConns(driverId int) []*websocket.Conn {
	value, ok := PassengerMap.Load(driverId)
	if !ok {
		return nil
	}
	return value.([]*websocket.Conn)
}

// Remove a connection
func Remove_PassConn(driverId int, conn *websocket.Conn) {
	value, ok := PassengerMap.Load(driverId)
	if !ok {
		return
	}
	conns := value.([]*websocket.Conn)
	for i, ws_conn := range conns {
		if ws_conn == conn {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	PassengerMap.Store(driverId, conns)
}
