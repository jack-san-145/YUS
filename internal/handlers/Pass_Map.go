package handlers

import (
	"fmt"
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

type SyncMapPassengerStore struct {
	PassMap sync.Map // driverID -> []*PassengerConn
}

func NewSyncMapPassengerStore() *SyncMapPassengerStore {
	return &SyncMapPassengerStore{}
}

// add new driver to the PassMap
func (s *SyncMapPassengerStore) AddDriver(driver_id int) {
	s.PassMap.Store(driver_id, []*PassengerConn{})
}

// remove driver from the PassMap
func (s *SyncMapPassengerStore) RemoveDriver(driver_id int) {
	s.PassMap.Delete(driver_id)
}

// add new passengers to the PassMap
func (s *SyncMapPassengerStore) AddPassengerConn(driverId int, conn *websocket.Conn) {
	value, _ := s.PassMap.Load(driverId)
	var conns []*PassengerConn
	if value != nil {
		conns = value.([]*PassengerConn)
	}
	conns = append(conns, &PassengerConn{Conn: conn}) //creating a new passenger connection with 'conn'
	s.PassMap.Store(driverId, conns)
}

// remove passenger connections from the PassMap
func (s *SyncMapPassengerStore) RemovePassengerConn(driverId int, conn *websocket.Conn) {
	value, ok := s.PassMap.Load(driverId)
	if !ok || value == nil {
		return
	}
	conns := value.([]*PassengerConn)
	newConns := make([]*PassengerConn, 0, len(conns))
	for _, c := range conns {
		if c.Conn != conn {
			newConns = append(newConns, c)
		}
	}
	s.PassMap.Store(driverId, newConns)
}

// Get passenger connection from the PassMap
func (s *SyncMapPassengerStore) GetPassengerConns(driverId int) []*PassengerConn {
	var conns []*PassengerConn
	value, ok := s.PassMap.Load(driverId)
	if !ok {
		return nil
	}

	if ok && value != nil { //to avoid the panic
		conns = value.([]*PassengerConn)
	} else {
		conns = []*PassengerConn{} // initialize a new slice
	}
	return conns
}

// send driver location updates to the passengers
func (s *SyncMapPassengerStore) BroadcastLocation(driver_id int, current_location models.Location) {

	passengers := s.GetPassengerConns(driver_id)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passengers))

	for _, p := range passengers {
		p.Mu.Lock() // serialize writes
		err := p.Conn.WriteJSON(current_location)
		p.Mu.Unlock()

		if err != nil {
			fmt.Println("error sending location to passenger:", err)
			p.Conn.Close() //closed the websocket connection
			s.RemovePassengerConn(driver_id, p.Conn)
		}
	}
}
