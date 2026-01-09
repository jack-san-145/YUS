package postgres

import (
	"context"
	"log"
	"yus/internal/models"
)

func (pg *PgStore) CheckRouteExistsForPassengerWS(ctx context.Context, route models.PassengerWsRequest) (bool, error) {
	var exists bool

	if route.DriverId == 1000 || route.RouteId == 0 { //1000 refers nil for driver and 0 refers nil for route in YUS
		return exists, nil
	}

	query := "select exists(Select 1 from current_bus_route where driver_id = $1 and route_id = $2 and direction = $3)"
	err := pg.Pool.QueryRow(ctx, query, route.DriverId, route.RouteId, route.Direction).Scan(&exists)
	if err != nil {
		log.Println("error while check the route exits for the passenger request - ", err)
		return exists, err
	}

	return exists, nil
}
