package passenger

import (
	"context"
	"fmt"
	"net/http"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

func (h *PassengerHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {

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

	go h.listenPassengerMessage(conn)

}

func (h *PassengerHandler) listenPassengerMessage(conn *websocket.Conn) {

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

			//used old_requested_bus_route.DriverId bcz the current requested_bus_route.DriverId not received , so we cleared the old connection
			PassengerConnStore.RemovePassengerConn(old_requested_bus_route.DriverId, conn)
			return
		}

		fmt.Println("requested_bus_route -", requested_bus_route)

		//to check if the passenger request ws route is present in current_bus_route,is only true when route_id,direction,driver_id matched
		if exists, _ := h.Store.DB.CheckRouteExistsForPassengerWS(context.Background(), requested_bus_route); exists {
			fmt.Printf("old driver - %v and new driver - %v\n", old_requested_bus_route.DriverId, requested_bus_route.DriverId)
			PassengerConnStore.RemovePassengerConn(old_requested_bus_route.DriverId, conn)

			//check if the driver exists on the passenger map
			ok := PassengerConnStore.DriverExists(requested_bus_route.DriverId)
			if !ok {

				//if driver doesn't exists add him to the passengerMap
				PassengerConnStore.AddDriver(requested_bus_route.DriverId)
			}

			//after that add passenger to the PassengerMap under that driver
			PassengerConnStore.AddPassengerConn(requested_bus_route.DriverId, conn)
			old_requested_bus_route = requested_bus_route
		} else {
			PassengerConnStore.RemovePassengerConn(old_requested_bus_route.DriverId, conn)
			old_requested_bus_route = requested_bus_route
		}
	}
}
