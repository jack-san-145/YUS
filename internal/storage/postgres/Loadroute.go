package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func Load_routes_by_src_and_dest(src string, dest string) {
	// var (
	// 	All_Routes []models.Route
	// )
	query := "select route_id from all_routes where (src = $1 and dest = $2) "
	all_routes_ids, err := pool.Query(context.Background(), query, src, dest)
	if err == sql.ErrNoRows {
		fmt.Print("error while accessing the all route id's - ", err)
		return
	} else if err != nil {
		fmt.Println("error while accessing the route id's - ", err)
		return
	}

	defer all_routes_ids.Close() // close the all_routes_ids pointer when this function ends

	for all_routes_ids.Next() {
		var (
			route models.Route
		)
		err := all_routes_ids.Scan(&route.Id)
		if err != nil {
			fmt.Println("error while scan the route_id - ", err)
			return
		}

		query = "select * from route_stops where route_id = $1"
		stops, err := pool.Query(context.Background(), query, route.Id)
		if err == sql.ErrNoRows {
			fmt.Println("empty row")
			continue
		} else if err != nil {
			fmt.Print("error while accessing the route_stops - ", err)
			continue
		}

		for stops.Next() {
			var route_stops models.RouteStops
			err := stops.Scan(
				&route.Id)
			fmt.Println(route_stops, err)

		}
		stops.Close()
	}

}

// function to load all up_routes
func Load_available_routes() []models.AvilableRoute {
	var Available_routes []models.AvilableRoute
	query := "select route_id,name,src,dest,direction from all_routes where direction = 'UP' "
	all_routes, err := pool.Query(context.Background(), query)
	if err != nil {
		fmt.Println("error while finding the the all_routes - ", err)
		return nil
	}

	defer all_routes.Close()

	for all_routes.Next() {
		var (
			bus_id int
			route  models.AvilableRoute
		)

		err := all_routes.Scan(&route.Id,
			&route.Name,
			&route.Src,
			&route.Dest,
			&route.Direction)
		if err != nil {
			fmt.Println("error while scanning route_id from all_routes - ", err)
			continue
		}
		route.Name = services.Convert_to_Normal(route.Name)

		query := "select bus_id from current_bus_route where route_id = $1"
		err = pool.QueryRow(context.Background(), query, route.Id).Scan(&bus_id)

		if errors.Is(err, sql.ErrNoRows) {
			//if route not present in bus_route and its available to map with a bus
			route.Available = true
		} else if err != nil {
			fmt.Println("error while scanning bus_id from current_bus_route - ", err)
			continue
		}

		if !route.Available {
			//if the bus present in the bus_route and its not available bcz it already mapped with a bus
			route.BusId = bus_id
		}
		Available_routes = append(Available_routes, route)
	}
	return Available_routes
}
