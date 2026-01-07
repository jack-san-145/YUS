package passenger

import (
	"fmt"
	"sync"
	"time"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

// struct that implements the PassConnStore
type MapPassengerStore struct {
	PassMap map[int][]*PassengerConn
	Rwm     sync.RWMutex
	Timers  map[int]*time.Timer
}

// method to return the MapPassengerStore
func NewMapPassengerStore() *MapPassengerStore {
	return &MapPassengerStore{
		PassMap: make(map[int][]*PassengerConn),
		Timers:  make(map[int]*time.Timer),
	}
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

	if _, exists := m.PassMap[driverID]; exists {
		m.CancelRemoval(driverID)
	} else {
		m.PassMap[driverID] = []*PassengerConn{}
	}

}

// remove driver from the PassMap
func (m *MapPassengerStore) RemoveDriver(driverID int) {

	m.Rwm.Lock()
	defer m.Rwm.Unlock()
	delete(m.PassMap, driverID)

}

// add new passengers to the PassMap
func (m *MapPassengerStore) AddPassengerConn(driverId int, conn *websocket.Conn) {

	m.RemovePassengerConn(driverId, conn) // to prevent concurrent passenger object creation

	m.Rwm.Lock()
	defer m.Rwm.Unlock()
	p := &PassengerConn{
		Conn: conn,
		Send: make(chan models.Location, 10), //buffer size 10 to store max 10 location updates
	}
	m.PassMap[driverId] = append(m.PassMap[driverId], p)
	go p.StartWriter(driverId, m) //start a new separate go routine to write location updates,so each passenger have their own writer which avoids blocking others
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

func (m *MapPassengerStore) TerminatePassengerConn(driverId int, conn *websocket.Conn) {
	m.Rwm.Lock()

	passengerconn_arr := m.PassMap[driverId]
	for _, p := range passengerconn_arr {
		if p.Conn == conn {
			p.CloseOnce.Do(func() {
				close(p.Send)
				p.Conn.Close()
			}) //closing the disconnected passenger channel
		}
	}
	m.Rwm.Unlock()
	m.RemovePassengerConn(driverId, conn)
}

// Get passenger connection from the PassMap
func (m *MapPassengerStore) GetPassengerConns(driverID int) []*PassengerConn {

	m.Rwm.RLock()
	defer m.Rwm.RUnlock()

	passengerConn := m.PassMap[driverID]
	copied := make([]*PassengerConn, len(passengerConn))
	copy(copied, passengerConn)

	return copied

}
func (p *PassengerConn) StartWriter(driverID int, m *MapPassengerStore) {
	const pingPeriod = 50 * time.Second
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {

		// 1 Location update
		case loc, ok := <-p.Send:
			if !ok {
				// channel closed → passenger removed
				return
			}

			p.Mu.Lock()
			err := p.Conn.WriteJSON(loc)
			p.Mu.Unlock()

			if err != nil {
				fmt.Println("error sending location to passenger:", err)
				p.Conn.Close()
				m.RemovePassengerConn(driverID, p.Conn)
				return
			}

		// 2️ Ping to keep connection alive
		case <-ticker.C:
			p.Mu.Lock()
			err := p.Conn.WriteMessage(websocket.PingMessage, nil)
			p.Mu.Unlock()

			if err != nil {
				fmt.Println("ping error:", err)
				p.Conn.Close()
				m.RemovePassengerConn(driverID, p.Conn)
				return
			}
		}
	}
}

// send driver location updates to the passengers
func (m *MapPassengerStore) BroadcastLocation(driver_id int, current_location models.Location) {

	passengers := m.GetPassengerConns(driver_id)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passengers))

	for _, p := range passengers {
		select {
		case p.Send <- current_location:
			//location update sents successfully
		default:
			//location update dropped
		}
	}
}

func (m *MapPassengerStore) ScheduleRemoval(driverID int) {
	m.Rwm.Lock()
	defer m.Rwm.Unlock()

	// if already scheduled, do nothing
	if _, exists := m.Timers[driverID]; exists {
		return
	}

	timer := time.AfterFunc(30*time.Minute, func() {

		m.Rwm.Lock()
		delete(m.PassMap, driverID)
		delete(m.Timers, driverID)
		m.Rwm.Unlock()
	})

	m.Timers[driverID] = timer
}

func (m *MapPassengerStore) CancelRemoval(driverID int) {
	m.Rwm.Lock()
	defer m.Rwm.Unlock()

	if timer, exists := m.Timers[driverID]; exists {
		timer.Stop() // stop timer
		delete(m.Timers, driverID)
	}
}
