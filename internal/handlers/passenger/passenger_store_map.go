package passenger

import (
	"fmt"
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

// struct that implements the PassConnStore
type MapPassengerStore struct {
	PassMap map[int][]*PassengerConn
	Rwm     sync.RWMutex
}

// method to return the MapPassengerStore
func NewMapPassengerStore() *MapPassengerStore {
	return &MapPassengerStore{PassMap: make(map[int][]*PassengerConn)}
}

// check if the driver exist or not
func (m *MapPassengerStore) DriverExists(driverID int) bool {

	m.Rwm.RLock()
	defer m.Rwm.RUnlock()
	_, ok := m.PassMap[driverID]
	return ok
}

// add new driver to the PassMap
func (m *MapPassengerStore) AddDriver(driverID int) {

	m.Rwm.Lock()
	defer m.Rwm.Unlock()
	m.PassMap[driverID] = []*PassengerConn{}

}

// remove driver from the PassMap
func (m *MapPassengerStore) RemoveDriver(driverID int) {

	m.Rwm.Lock()
	defer m.Rwm.Unlock()
	delete(m.PassMap, driverID)

}

// add new passengers to the PassMap
func (m *MapPassengerStore) AddPassengerConn(driverId int, conn *websocket.Conn) {

	m.Rwm.Lock()
	defer m.Rwm.Unlock()
	m.PassMap[driverId] = append(m.PassMap[driverId], &PassengerConn{Conn: conn})

}

// remove passenger connections from the PassMap
func (m *MapPassengerStore) RemovePassengerConn(driverId int, conn *websocket.Conn) {

	m.Rwm.Lock()
	defer m.Rwm.Unlock()

	passengerconn_arr := m.PassMap[driverId]

	for i, p := range passengerconn_arr {
		if p.Conn == conn {
			passengerconn_arr = append(passengerconn_arr[:i], passengerconn_arr[i+1:]...)
			m.PassMap[driverId] = passengerconn_arr
			return
		}
	}

}

// Get passenger connection from the PassMap
func (m *MapPassengerStore) GetPassengerConns(driverID int) []*PassengerConn {

	m.Rwm.RLock()
	defer m.Rwm.RUnlock()

	passangerConn := m.PassMap[driverID]
	copied := make([]*PassengerConn, len(passangerConn))
	copy(copied, passangerConn)

	return copied

}

// send driver location updates to the passengers
func (m *MapPassengerStore) BroadcastLocation(driver_id int, current_location models.Location) {

	passengers := m.GetPassengerConns(driver_id)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passengers))

	for _, p := range passengers {
		p.Mu.Lock() // serialize writes
		err := p.Conn.WriteJSON(current_location)
		p.Mu.Unlock()

		if err != nil {
			fmt.Println("error sending location to passenger:", err)
			p.Conn.Close() //closed the websocket connection
			m.RemovePassengerConn(driver_id, p.Conn)
		}

	}
}
