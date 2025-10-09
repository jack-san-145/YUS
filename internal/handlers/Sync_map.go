package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
)

var ConnMap sync.Map // key: int, value: []*websocket.Conn (

// key = driver_id and value = passenger Websocket conncetions

// Store a connection
func AddConn(driverId int, conn *websocket.Conn) {
	value, _ := ConnMap.LoadOrStore(driverId, []*websocket.Conn{})
	conns := value.([]*websocket.Conn) //assings the passenger ws connections to the conns
	conns = append(conns, conn)        //add the new passenger to the existing array and it mapped with driverId
	ConnMap.Store(driverId, conns)     //store the final conns to the ConnMap
}

// Get connections
func GetConns(driverId int) []*websocket.Conn {
	value, ok := ConnMap.Load(driverId)
	if !ok {
		return nil
	}
	return value.([]*websocket.Conn)
}

// Remove a connection
func RemoveConn(driverId int, conn *websocket.Conn) {
	value, ok := ConnMap.Load(driverId)
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
	ConnMap.Store(driverId, conns)
}
