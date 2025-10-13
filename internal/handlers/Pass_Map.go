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

// package handlers

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// 	"yus/internal/models"

// 	"github.com/gorilla/websocket"
// )

// var PassengerMap sync.Map // key: int, value: []*websocket.Conn

// // key = driver_id & value = passenger Websocket conncetions

// // Store a connection
// func Add_PassConn(driverId int, conn *websocket.Conn) {
// 	value, ok := PassengerMap.Load(driverId)

// 	// conns := value.([]*websocket.Conn) //assings the passenger ws connections to the conns
// 	// conns = append(conns, conn)        //add the new passenger to the existing array and it mapped with driverId
// 	// fmt.Println("after added the passenger ws conn - ", conn)

// 	var conns []*websocket.Conn
// 	if ok && value != nil { //to avoid the panic
// 		conns = value.([]*websocket.Conn)
// 	} else {
// 		fmt.Println("slice is nil")
// 		conns = []*websocket.Conn{} // initialize a new slice
// 	}
// 	conns = append(conns, conn)         // add new passenger
// 	PassengerMap.Store(driverId, conns) //store the final conns to the ConnMap
// }

// // Remove a connection
// func Remove_PassConn(driverId int, conn *websocket.Conn) {
// 	if driverId != 0 {
// 		var conns []*websocket.Conn
// 		value, ok := PassengerMap.Load(driverId)
// 		if !ok {
// 			return
// 		}

// 		if ok && value != nil {
// 			conns = value.([]*websocket.Conn) // safe now
// 		} else {
// 			conns = []*websocket.Conn{} // initialize a new slice
// 		}
// 		for i, ws_conn := range conns {
// 			if ws_conn == conn {
// 				conns = append(conns[:i], conns[i+1:]...) //appending all the other passenger connections without the specified ones
// 				break
// 			}
// 		}

// 		PassengerMap.Store(driverId, conns) //store the remaining connections to the passengerMap
// 	}
// }

// // func Send_location_to_passenger(driver_id int, current_location models.Location) {
// // 	value, ok := PassengerMap.Load(driver_id)
// // 	if !ok || value == nil {
// // 		fmt.Printf("driver_id - %v sending location: 0 passengers connected\n", driver_id)
// // 		return
// // 	}

// // 	passenger_conns := value.([]*websocket.Conn)
// // 	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passenger_conns))

// // 	for _, pass_conn := range passenger_conns {
// // 		go func(conn *websocket.Conn) {
// // 			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
// // 			err := conn.WriteJSON(current_location)
// // 			if err != nil {
// // 				fmt.Println("error sending location to passenger:", err)
// // 				conn.Close()
// // 				Remove_PassConn(driver_id, conn)
// // 			}
// // 		}(pass_conn)
// // 	}
// // }

// // func Send_location_to_passenger(driver_id int, current_location models.Location) {

// // 	var conn_arr []*websocket.Conn
// // 	value, ok := PassengerMap.Load(driver_id)
// // 	fmt.Println("")
// // 	fmt.Printf("driver_id - %v sending location", driver_id)

// // 	if ok && value != nil {
// // 		conn_arr = value.([]*websocket.Conn) // to avoid panic
// // 	} else {
// // 		conn_arr = []*websocket.Conn{} // initialize a new slice
// // 	}

// // 	if ok {
// // 		passenger_conns_for_driverID := conn_arr //passenger ws under specific driver_id

// // 		fmt.Println("")
// // 		fmt.Printf(" users = %v connected driver = %v -", len(passenger_conns_for_driverID), driver_id)
// // 		for _, pass_conn := range passenger_conns_for_driverID {
// // 			err := pass_conn.WriteJSON(current_location)
// // 			if err != nil {
// // 				fmt.Println("error while sending the current_location to passenger - ", err)
// // 				Remove_PassConn(driver_id, pass_conn)

// // 			}
// // 		}

// // 	}
// // }
