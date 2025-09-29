package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"yus/internal/models"
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
