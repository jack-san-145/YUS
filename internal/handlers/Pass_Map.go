package handlers

import (
	"fmt"
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

// to allow only one go routine to write data in one ws conn at a time
type PassengerConn struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

var PassengerMap sync.Map // driverId -> []*PassengerConn

func Add_PassConn(driverId int, conn *websocket.Conn) {
	value, _ := PassengerMap.Load(driverId)
	var conns []*PassengerConn
	if value != nil {
		conns = value.([]*PassengerConn)
	}
	conns = append(conns, &PassengerConn{Conn: conn}) //creating a new passenger connection with 'conn'
	PassengerMap.Store(driverId, conns)
}

func Remove_PassConn(driverId int, conn *websocket.Conn) {
	value, ok := PassengerMap.Load(driverId)
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
	PassengerMap.Store(driverId, newConns)
}

func Send_location_to_passenger(driver_id int, current_location models.Location) {
	value, ok := PassengerMap.Load(driver_id)
	if !ok || value == nil {
		fmt.Printf("driver_id - %v sending location: 0 passengers connected\n", driver_id)
		return
	}
	passengers := value.([]*PassengerConn)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passengers))

	for _, p := range passengers {
		p.Mu.Lock() // serialize writes
		err := p.Conn.WriteJSON(current_location)
		p.Mu.Unlock()

		if err != nil {
			fmt.Println("error sending location to passenger:", err)
			p.Conn.Close() //closed the websocket connection
			Remove_PassConn(driver_id, p.Conn)
		}
	}
}

func Remove_Driver_from_passengerMap(driver_id int) {
	PassengerMap.Delete(driver_id)
}

func Add_Driver_to_passengerMap(driver_id int) {
	PassengerMap.Store(driver_id, []*PassengerConn{})
}

// Get connections
func Get_PassConns(driverId int) []*PassengerConn {
	var conns []*PassengerConn
	value, ok := PassengerMap.Load(driverId)
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
