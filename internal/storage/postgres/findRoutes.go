package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
)

func FindRoutes_by_src_dest(src string, dest string) []models.CurrentRoute {
	var All_routes []models.CurrentRoute
	query := "select bus_id,driver_id,route_id,direction,name,src,dest from current_bus_route where src = $1 and dest = $2"
	routes, err := pool.Query(context.Background(), query, src, dest)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("No routes available for this src and dest")
	} else if err != nil {
		fmt.Println("error while select the routes by src and dest from current_bus_route - ", err)
	}

	defer routes.Close()
	for routes.Next() {
		var route models.CurrentRoute
		routes.Scan(&route.BusId,
			&route.DriverId,
			&route.RouteId,
			&route.Direction,
			&route.RouteName,
			&route.Src,
			&route.Dest)

		All_routes = append(All_routes, route)
	}
	return All_routes
}
