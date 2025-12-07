package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Passenger_Conn struct {
	conn *websocket.Conn
	ru   sync.RWMutex
}

var common_ru sync.RWMutex

var Passenger_Map = make(map[int][]Passenger_Conn)

func Add_Driver_to_passenger_Map(driverID int) {

	common_ru.Lock()
	Passenger_Map[driverID] = []Passenger_Conn{}
	common_ru.Unlock()
}

func Remove_Driver_from_passenger_Map(driverID int) {

	common_ru.Lock()
	delete(Passenger_Map, driverID)
	common_ru.Unlock()
}

func Get_Pass_Conns(driverID int) []Passenger_Conn {

	common_ru.RLock()
	passangerconn := Passenger_Map[driverID]
	common_ru.RUnlock()

	return passangerconn
}
