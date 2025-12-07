package handlers

import (
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

// to allow only one go routine to write data in one ws conn at a time
type PassengerConn struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

// interface to implement passenger map operations
type PassengerConnStore interface {
	AddDriver(driverID int)
	RemoveDriver(driverID int)
	AddPassengerConn(driverID int, conn *websocket.Conn)
	RemovePassengerConn(driverID int, conn *websocket.Conn)
	GetPassengerConns(driverID int) []*PassengerConn
	BroadcastLocation(driverID int, location models.Location)
}
