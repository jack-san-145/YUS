package handlers

import (
	"fmt"
	"sync"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

type Passenger_Conn struct {
	Conn *websocket.Conn
	Rwm  sync.RWMutex
}

var common_ru sync.RWMutex

var Passenger_Map = make(map[int][]*Passenger_Conn)

func Add_Driver_to_passenger_Map(driverID int) {

	common_ru.Lock()
	Passenger_Map[driverID] = []*Passenger_Conn{}
	common_ru.Unlock()
}

func Remove_Driver_from_passenger_Map(driverID int) {

	common_ru.Lock()
	delete(Passenger_Map, driverID)
	common_ru.Unlock()
}

func Get_Pass_Conns(driverID int) []*Passenger_Conn {

	common_ru.RLock()
	passangerconn := Passenger_Map[driverID]
	common_ru.RUnlock()

	return passangerconn
}

func Add_Pass_Conn(driverId int, conn *websocket.Conn) {

	common_ru.Lock()
	Passenger_Map[driverId] = append(Passenger_Map[driverId], &Passenger_Conn{Conn: conn})
	common_ru.Unlock()
}

func Remove_Pass_Conn(driverId int, conn *websocket.Conn) {

	common_ru.Lock()

	defer common_ru.Unlock()

	passengerconn_arr := Passenger_Map[driverId]

	for i, p := range passengerconn_arr {
		if p.Conn == conn {
			passengerconn_arr = append(passengerconn_arr[:i], passengerconn_arr[i+1:]...)
			Passenger_Map[driverId] = passengerconn_arr
			return
		}
	}

}

func Send_location_to_passengers(driver_id int, current_location models.Location) {

	passengers := Get_Pass_Conns(driver_id)
	fmt.Printf("driver_id - %v sending location, users = %d\n", driver_id, len(passengers))

	for _, p := range passengers {
		p.Rwm.Lock() // serialize writes
		err := p.Conn.WriteJSON(current_location)
		p.Rwm.Unlock()

		if err != nil {
			fmt.Println("error sending location to passenger:", err)
			p.Conn.Close() //closed the websocket connection
			Remove_Pass_Conn(driver_id, p.Conn)
		}
	}
}
