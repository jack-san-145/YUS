package postgres

import (
	"context"
	"fmt"
	"yus/internal/models"
)

func (pg *PgStore) CheckRouteExistsForPassengerWS(ctx context.Context, route models.PassengerWsRequest) (bool, error) {
	var exists bool

	fmt.Printf("route.Direction - %v & type - %T ", route.Direction, route.Direction)
	fmt.Printf("route.DriverId - %v & type - %T ", route.DriverId, route.DriverId)
	fmt.Printf("route.RouteId - %v & type - %T ", route.RouteId, route.RouteId)
	query := "select exists(Select 1 from current_bus_route where driver_id = $1 and route_id = $2 and direction = $3)"
	err := pg.Pool.QueryRow(ctx, query, route.DriverId, route.RouteId, route.Direction).Scan(&exists)
	if err != nil {
		fmt.Println("error while check the route exits for the passenger request - ", err)
		return exists, err
	}
	fmt.Println("exists - ", exists)

	return exists, nil
}
