package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
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

		route.RouteName = services.Convert_to_Normal(route.RouteName)
		route.Src = services.Convert_to_Normal(route.Src)
		route.Dest = services.Convert_to_Normal(route.Dest)

		route.Stops = findStops(route.RouteId, route.Direction)
		All_routes = append(All_routes, route)
	}
	return All_routes
}

func findStops(route_id int, direction string) []models.RouteStops {
	var route_stops []models.RouteStops
	query := "select stop_name,is_stop,lat,lon,arrival_time,departure_time from route_stops where route_id = $1 and direction = $2"
	all_stops, err := pool.Query(context.Background(), query, route_id, direction)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("no stops found for this route_id ")
	} else if err != nil {
		fmt.Println("error while finding the route stops - ", err)
	}

	defer all_stops.Close()
	for all_stops.Next() {
		var stop models.RouteStops
		all_stops.Scan(&stop.LocationName,
			&stop.IsStop,
			&stop.Lat,
			&stop.Lon,
			&stop.Arrival_time,
			&stop.Departure_time)

		stop.LocationName = services.Convert_to_Normal(stop.LocationName)
		route_stops = append(route_stops, stop)
	}
	return route_stops
}
