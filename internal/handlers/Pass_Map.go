package handlers

import (
	"fmt"
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

var PassengerMap sync.Map // key: int, value: []*websocket.Conn

// key = driver_id & value = passenger Websocket conncetions

// Store a connection
func Add_PassConn(driverId int, conn *websocket.Conn) {
	value, ok := PassengerMap.Load(driverId)

	// conns := value.([]*websocket.Conn) //assings the passenger ws connections to the conns
	// conns = append(conns, conn)        //add the new passenger to the existing array and it mapped with driverId
	// fmt.Println("after added the passenger ws conn - ", conn)

	var conns []*websocket.Conn
	if ok && value != nil { //to avoid the panic
		conns = value.([]*websocket.Conn)
	} else {
		conns = []*websocket.Conn{} // initialize a new slice
	}
	PassengerMap.Store(driverId, conns) //store the final conns to the ConnMap
}

func Add_Driver_to_passengerMap(driver_id int) {
	PassengerMap.Store(driver_id, []*websocket.Conn{})
}

func Remove_Driver_from_passengerMap(driver_id int) {
	PassengerMap.Delete(driver_id)
}

// Get connections
func Get_PassConns(driverId int) []*websocket.Conn {
	var conns []*websocket.Conn
	value, ok := PassengerMap.Load(driverId)
	if !ok {
		return nil
	}

	if ok && value != nil { //to avoid the panic
		conns = value.([]*websocket.Conn)
	} else {
		conns = []*websocket.Conn{} // initialize a new slice
	}
	return conns
}

// Remove a connection
func Remove_PassConn(driverId int, conn *websocket.Conn) {
	if driverId != 0 {
		var conns []*websocket.Conn
		value, ok := PassengerMap.Load(driverId)
		if !ok {
			return
		}

		if ok && value != nil {
			conns = value.([]*websocket.Conn) // safe now
		} else {
			conns = []*websocket.Conn{} // initialize a new slice
		}
		for i, ws_conn := range conns {
			if ws_conn == conn {
				conns = append(conns[:i], conns[i+1:]...) //appending all the other passenger connections without the specified ones
				break
			}
		}

		PassengerMap.Store(driverId, conns) //store the remaining connections to the passengerMap
	}
}

func Send_location_to_passenger(driver_id int, current_location models.Location) {

	var conn_arr []*websocket.Conn
	value, ok := PassengerMap.Load(driver_id)
	fmt.Printf("driver_id - %v & ok - %v ", driver_id, ok)

	if ok && value != nil {
		conn_arr = value.([]*websocket.Conn) // to avoid panic
	} else {
		conn_arr = []*websocket.Conn{} // initialize a new slice
	}

	if ok {
		passenger_conns_for_driverID := conn_arr //passenger ws under specific driver_id

		fmt.Printf("no of users = %v connected with the driver = %v -", len(passenger_conns_for_driverID), driver_id)
		for _, pass_conn := range passenger_conns_for_driverID {
			err := pass_conn.WriteJSON(current_location)
			if err != nil {
				Remove_PassConn(driver_id, pass_conn)
				fmt.Println("error while sending the current_location to passenger - ", err)
			}
		}

	}
}
