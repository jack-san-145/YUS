package passenger

import (
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

// to allow only one go routine to write data in one ws conn at a time
type PassengerConn struct {
	Conn      *websocket.Conn
	Mu        sync.Mutex
	Send      chan models.Location
	CloseOnce sync.Once
}

// interface to implement passenger map operations
type PassengerConnectionManager interface {
	DriverExists(driverID int) bool
	AddDriver(driverID int)
	RemoveDriver(driverID int)
	AddPassengerConn(driverID int, conn *websocket.Conn)
	RemovePassengerConn(driverID int, conn *websocket.Conn)
	GetPassengerConns(driverID int) []*PassengerConn
	BroadcastLocation(driverID int, location models.Location)
	ScheduleRemoval(driverID int)
	CancelRemoval(driverID int)
}

var PassengerConnStore = NewMapPassengerStore() //  to store the passenger connections in normal Go Map

// var PassengerConnStore = NewSyncMapPassengerStore() // to store the passenger connections in Sync.Map
