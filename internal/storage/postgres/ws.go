package postgres

import (
	"context"
	"fmt"
	"yus/internal/models"
)

func Check_route_exits_for_pass_Ws(route models.PassengerWsRequest) bool {
	var exists bool
	query := "select exists(Select 1 from current_bus_route where driver_id = $1 and route_id = $2 and direction = $3)"
	err := pool.QueryRow(context.Background(), query, route.DriverId, route.RouteId, route.Direction).Scan(&exists)
	if err != nil {
		fmt.Println("error while check the route exits for the passenger request - ", err)
	}

	return exists
}
