package passenger

import (
	"context"
	// "log"
	"net/http"
	"yus/internal/models"

	"github.com/gorilla/websocket"
)

func (h *PassengerHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

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
			} else {
			}

			//used old_requested_bus_route.DriverId bcz the current requested_bus_route.DriverId not received , so we cleared the old connection
			PassengerConnStore.TerminatePassengerConn(old_requested_bus_route.DriverId, conn)

			return
		}

		//to check if the passenger request ws route is present in current_bus_route,is only true when route_id,direction,driver_id matched
		if exists, _ := h.Store.DB.CheckRouteExistsForPassengerWS(context.Background(), requested_bus_route); exists {
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
