package handlers

import (
	"fmt"
	"net/http"
	"time"
	"yus/internal/models"
	"yus/internal/storage/postgres"

	"github.com/gorilla/websocket"
)

// AddDriver(driverID int)
// 	RemoveDriver(driverID int)
// 	AddPassengerConn(driverID int, conn *websocket.Conn)
// 	RemovePassengerConn(driverID int, conn *websocket.Conn)
// 	GetPassengerConns(driverID int) []*PassengerConn
// 	BroadcastLocation(driverID int, location models.Location)

func Passenger_Ws_handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("working passenger app websocket upgradation")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}
	fmt.Println("passenger connected successfully ")

	go listen_passenger_message(conn)
	// Passenger_All_WS_connections = append(Passenger_All_WS_connections, conn)

}

func listen_passenger_message(conn *websocket.Conn) {
	const pingPeriod = 50 * time.Second

	// Ping ticker to keep connection alive
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Println("ping error:", err)
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()

	var (
		old_requested_bus_route models.PassengerWsRequest
		requested_bus_route     models.PassengerWsRequest
	)

	for {
		// Read passenger requests asynchronously
		err := conn.ReadJSON(&requested_bus_route)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("passenger connection closed unexpectedly -", err)
			} else {
				fmt.Println("error reading passenger ws message -", err)
			}
			close(done)

			//used old_requested_bus_route.DriverId bcz the current requested_bus_route.DriverId not received , so we cleared the old connection
			Remove_PassConn(old_requested_bus_route.DriverId, conn)
			return
		}

		fmt.Println("requested_bus_route -", requested_bus_route)

		//to check if the passenger request ws route is present in current_bus_route,is only true when route_id,direction,driver_id matched
		if postgres.Check_route_exits_for_pass_Ws(requested_bus_route) {
			fmt.Printf("old driver - %v and new driver - %v\n", old_requested_bus_route.DriverId, requested_bus_route.DriverId)
			Remove_PassConn(old_requested_bus_route.DriverId, conn)

			//check if the driver exists on the passenger map
			_, ok := PassengerMap.Load(requested_bus_route.DriverId)
			if !ok {

				//if driver doesn't exists add him to the passengerMap
				Add_Driver_to_passengerMap(requested_bus_route.DriverId)
			}

			//after that add passenger to the PassengerMap under that driver
			Add_PassConn(requested_bus_route.DriverId, conn)
			old_requested_bus_route = requested_bus_route
		} else {
			Remove_PassConn(old_requested_bus_route.DriverId, conn)
			old_requested_bus_route = requested_bus_route
		}
	}
}
