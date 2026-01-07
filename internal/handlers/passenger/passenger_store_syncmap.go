package passenger

import (
	"fmt"
	"sync"
	"time"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

type SyncMapPassengerStore struct {
	PassMap sync.Map // driverID -> []*PassengerConn
	Rwm     sync.RWMutex
	Timers  map[int]*time.Timer
}

func NewSyncMapPassengerStore() *SyncMapPassengerStore {
	return &SyncMapPassengerStore{}
}

// check if the driver exist or not
func (s *SyncMapPassengerStore) DriverExists(driverID int) bool {
	_, ok := s.PassMap.Load(driverID)
	return ok
}

// add new driver to the PassMap
func (s *SyncMapPassengerStore) AddDriver(driverID int) {
	s.PassMap.Store(driverID, []*PassengerConn{})
}

// remove driver from the PassMap
func (s *SyncMapPassengerStore) RemoveDriver(driverID int) {
	s.PassMap.Delete(driverID)
}

// add new passengers to the PassMap
func (s *SyncMapPassengerStore) AddPassengerConn(driverID int, conn *websocket.Conn) {
	s.RemovePassengerConn(driverID, conn) // prevent duplicate connections

	value, _ := s.PassMap.Load(driverID)
	var conns []*PassengerConn
	if value != nil {
		conns = value.([]*PassengerConn)
	}

	p := &PassengerConn{
		Conn: conn,
		Send: make(chan models.Location, 10), // buffered channel
	}
	conns = append(conns, p)
	s.PassMap.Store(driverID, conns)

	go p.StartWriterSMP(driverID, s)
}

// remove passenger connections from the PassMap
func (s *SyncMapPassengerStore) RemovePassengerConn(driverID int, conn *websocket.Conn) {
	value, ok := s.PassMap.Load(driverID)
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
	s.PassMap.Store(driverID, newConns)
}

func (s *SyncMapPassengerStore) TerminatePassengerConn(driverId int, conn *websocket.Conn) {

	passengerconn_arr := s.GetPassengerConns(driverId)
	for _, p := range passengerconn_arr {
		if p.Conn == conn {
			p.CloseOnce.Do(func() {
				close(p.Send)
				p.Conn.Close()
			}) //closing the disconnected passenger channel
		}
	}
	s.RemovePassengerConn(driverId, conn)

}

// Get passenger connections
func (s *SyncMapPassengerStore) GetPassengerConns(driverID int) []*PassengerConn {
	value, ok := s.PassMap.Load(driverID)
	if !ok || value == nil {
		return nil
	}

	conns := value.([]*PassengerConn)
	copied := make([]*PassengerConn, len(conns))
	copy(copied, conns)
	return copied
}

// Start writer goroutine for a passenger
func (p *PassengerConn) StartWriterSMP(driverID int, s *SyncMapPassengerStore) {
	const pingPeriod = 50 * time.Second
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		// 1. Location update
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
				p.CloseOnce.Do(func() {
					p.Conn.Close()
				})
				s.RemovePassengerConn(driverID, p.Conn)
				return
			}

		// 2. Ping to keep connection alive
		case <-ticker.C:
			p.Mu.Lock()
			err := p.Conn.WriteMessage(websocket.PingMessage, nil)
			p.Mu.Unlock()

			if err != nil {
				fmt.Println("ping error:", err)
				p.CloseOnce.Do(func() {
					p.Conn.Close()
				})
				s.RemovePassengerConn(driverID, p.Conn)
				return
			}
		}
	}
}

// Broadcast driver location to all passengers
func (s *SyncMapPassengerStore) BroadcastLocation(driverID int, loc models.Location) {
	passengers := s.GetPassengerConns(driverID)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driverID, len(passengers))

	for _, p := range passengers {
		select {
		case p.Send <- loc:
			// sent successfully
		default:
			// dropped if channel full
		}
	}
}

func (s *SyncMapPassengerStore) ScheduleRemoval(driverID int) {
	s.Rwm.Lock()
	defer s.Rwm.Unlock()

	// if already scheduled, do nothing
	if _, exists := s.Timers[driverID]; exists {
		return
	}

	timer := time.AfterFunc(30*time.Minute, func() {

		s.Rwm.Lock()
		s.PassMap.Delete(driverID)
		delete(s.Timers, driverID)
		s.Rwm.Unlock()
	})

	s.Timers[driverID] = timer
}

func (s *SyncMapPassengerStore) CancelRemoval(driverID int) {
	s.Rwm.Lock()
	defer s.Rwm.Unlock()

	if timer, exists := s.Timers[driverID]; exists {
		timer.Stop() // stop timer
		delete(s.Timers, driverID)
	}
}
