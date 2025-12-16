package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func FindRouteByBusOrDriverID(ctx context.Context, busID int, requestFrom string) (*models.AllRoute, error) {
	if requestFrom == "DRIVER" {
		driver_id := &busID
		query := "select bus_id from current_bus_route where driver_id = $1"
		err := pool.QueryRow(context.Background(), query, driver_id).Scan(&busID)
		if err != nil {
			fmt.Println("error while finding bus_id by driver_id - ", err)
		}
	}

	var r = &models.AllRoute{}

	query := "select bus_id,driver_id,route_id,direction,route_name,src,dest from current_bus_route where bus_id = $1 "
	err := pool.QueryRow(context.Background(), query, busID).Scan(&r.Currentroute.BusId,
		&r.Currentroute.DriverId,
		&r.Currentroute.RouteId,
		&r.Currentroute.Direction,
		&r.Currentroute.RouteName,
		&r.Currentroute.Src,
		&r.Currentroute.Dest)

	r.Currentroute.RouteName = services.Convert_to_Normal(r.Currentroute.RouteName)
	r.Currentroute.Src = services.Convert_to_Normal(r.Currentroute.Src)
	r.Currentroute.Dest = services.Convert_to_Normal(r.Currentroute.Dest)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("no route found for this bus_id ")
		return r, fmt.Errorf("no route found")
	} else if err != nil {
		fmt.Println("error while finding the route with bus_id - ", err)
		return r, err
	}

	FindStops(ctx, &r.Currentroute)
	if r.Currentroute.Direction == "UP" {
		r.Uproute = r.Currentroute
		r.Uproute.Active = true
		FindStops(ctx, &r.Uproute)

		r.Currentroute.Direction = "DOWN"
		r.Downroute = r.Currentroute
		r.Downroute.Src, r.Downroute.Dest = r.Downroute.Dest, r.Downroute.Src

		r.Downroute.Active = false
		FindStops(ctx, &r.Downroute)

	} else if r.Currentroute.Direction == "DOWN" {
		r.Downroute = r.Currentroute
		r.Downroute.Active = true
		FindStops(ctx, &r.Downroute)

		r.Currentroute.Direction = "UP"
		r.Uproute = r.Currentroute
		r.Uproute.Src, r.Uproute.Dest = r.Uproute.Dest, r.Uproute.Src

		r.Uproute.Active = false
		FindStops(ctx, &r.Uproute)
	}

	return r, nil
	//here route is an current route
}

func FindRoutesBySrcDst(ctx context.Context, src string, dest string) ([]models.CurrentRoute, error) {
	var (
		All_routes []models.CurrentRoute
		found      bool
	)
	query := "select bus_id,driver_id,route_id,direction,route_name,src,dest from current_bus_route where src = $1 and dest = $2"
	routes, err := pool.Query(context.Background(), query, src, dest)
	if err != nil {
		fmt.Println("error while select the routes by src and dest from current_bus_route - ", err)
		return All_routes, err
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

		FindStops(ctx, &route) //it sets the stops to the route with pointer
		All_routes = append(All_routes, route)
		found = true
	}

	if found {
		fmt.Println("succesfully route founded")
	} else if !routes.Next() {
		fmt.Println("No routes available for this src and dest,so find reverse route")
		All_routes, _ = FindReverseRoutesBySrcDest(ctx, dest, src)

	}
	return All_routes, nil

}

func FindStops(ctx context.Context, route *models.CurrentRoute) error {
	var (
		route_stops []models.RouteStops
		route_name  string
	)
	query := `
			select route_name,stop_sequence, stop_name, is_stop, lat, lon, arrival_time, departure_time
			from route_stops
			where route_id = $1 AND direction = $2
			order BY stop_sequence asc
		`

	all_stops, err := pool.Query(context.Background(), query, route.RouteId, route.Direction)
	if err != nil {
		fmt.Println("error while finding the route stops - ", err)
		return err
	}

	defer all_stops.Close()
	for all_stops.Next() {
		var (
			stop models.RouteStops
		)
		all_stops.Scan(&route_name,
			&stop.StopSequence,
			&stop.LocationName,
			&stop.IsStop,
			&stop.Lat,
			&stop.Lon,
			&stop.Arrival_time,
			&stop.Departure_time)

		stop.LocationName = services.Convert_to_Normal(stop.LocationName)
		route_stops = append(route_stops, stop)
	}
	route.RouteName = services.Convert_to_Normal(route_name)
	route.Stops = route_stops
	return nil
}

func FindReverseRoutesBySrcDest(ctx context.Context, src string, dest string) ([]models.CurrentRoute, error) {

	var (
		All_routes []models.CurrentRoute
		found      bool
	)

	query := "select bus_id,driver_id,route_id,direction,route_name,src,dest from current_bus_route where src = $1 and dest = $2"
	routes, err := pool.Query(context.Background(), query, src, dest)
	if err != nil {
		fmt.Println("error while select the reverse routes by src and dest from current_bus_route - ", err)
		return All_routes, err
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

		route.Src = dest //replace the original src given by the passenger
		route.Dest = src //replace the original dest given by the passenger

		if route.Direction == "UP" {
			route.Direction = "DOWN"
		} else if route.Direction == "DOWN" {
			route.Direction = "UP"
		}

		FindStops(ctx, &route) //it sets the stops to the route with pointer
		All_routes = append(All_routes, route)
		found = true
	}

	if found {
		fmt.Println("succesfully finds reverse route ")
	} else if !routes.Next() {
		fmt.Println("No reverse routes available for this src and dest")
	}
	return All_routes, nil

}

/*
1. check the direction exists on db
2. if exists check all the current routes with routestops
3. if doesn't exists go for all_routes and then find the match with the routestops
4. return the matched routes
*/
func FindRoutesBySrcDstStop(ctx context.Context, original_src string, original_dest string, original_stop string) ([]models.CurrentRoute, error) {

	var (
		Matched_routes   []models.CurrentRoute
		temp_src         string
		temp_dest        string
		filterWith       string
		direction        string
		direction_exists bool
	)

	if original_src == "Kcet" {
		direction = "DOWN"
		temp_src = original_src
		temp_dest = original_stop //changed
		filterWith = "dest"
		//find is stop for dest
	} else if original_dest == "Kcet" {
		direction = "UP"
		temp_dest = original_dest
		temp_src = original_stop //changed
		filterWith = "src"
		//find is stop for src
	}
	fmt.Println("filter with - ", filterWith)

	query := "select exists(Select 1 from current_bus_route where direction = $1) "
	err := pool.QueryRow(context.Background(), query, direction).Scan(&direction_exists)
	if err != nil {
		fmt.Println("error while check the existance of direction - ", err)
		return Matched_routes, err
	}

	if direction_exists {
		query = `SELECT c.bus_id, c.driver_id, c.route_id, c.direction, c.route_name, c.src, c.dest,
					CASE 
						WHEN c.direction = 'UP' THEN rs_src.is_stop
						WHEN c.direction = 'DOWN' THEN rs_dest.is_stop
					END AS is_stop
				FROM current_bus_route c
				JOIN route_stops rs_src
					ON rs_src.route_id = c.route_id
					AND rs_src.direction = c.direction
				JOIN route_stops rs_dest
					ON rs_dest.route_id = c.route_id
					AND rs_dest.direction = c.direction
				WHERE rs_src.stop_name LIKE $1
				AND rs_dest.stop_name LIKE $2
				AND rs_src.stop_sequence < rs_dest.stop_sequence
				ORDER BY
				  CASE 
						WHEN c.direction = 'UP' THEN rs_src.is_stop
						WHEN c.direction = 'DOWN' THEN rs_dest.is_stop
				  END DESC,
				rs_src.stop_sequence;`

		rows, err := pool.Query(context.Background(), query, temp_src+"%", temp_dest+"%")
		if err != nil {
			fmt.Println("error while finding the active route which present in current_bus_route - ", err)
		}

		defer rows.Close()

		for rows.Next() {
			var route models.CurrentRoute
			rows.Scan(&route.BusId,
				&route.DriverId,
				&route.RouteId,
				&route.Direction,
				&route.RouteName,
				&route.Src,
				&route.Dest,
				&route.IsStop)

			FindStops(ctx, &route)
			fmt.Println("active route - ", route)
			Matched_routes = append(Matched_routes, route)
		}

	} else {
		query = `SELECT 
				c.bus_id,
				c.driver_id,
				c.route_id,
				CASE WHEN c.direction = 'UP' THEN 'DOWN' ELSE 'UP' END AS direction,
				c.route_name,
				c.src,
				c.dest,
				CASE 
					WHEN c.direction = 'UP' THEN rs_dest.is_stop   -- swapped
					WHEN c.direction = 'DOWN' THEN rs_src.is_stop  -- swapped
				END AS is_stop
			FROM current_bus_route c
			JOIN route_stops rs_src
				ON rs_src.route_id = c.route_id
				AND rs_src.direction = CASE WHEN c.direction = 'UP' THEN 'DOWN' ELSE 'UP' END
			JOIN route_stops rs_dest
				ON rs_dest.route_id = c.route_id
				AND rs_dest.direction = CASE WHEN c.direction = 'UP' THEN 'DOWN' ELSE 'UP' END
			WHERE rs_src.stop_name LIKE $1
			AND rs_dest.stop_name LIKE $2
			AND rs_src.stop_sequence < rs_dest.stop_sequence
			ORDER BY
				CASE 
					WHEN c.direction = 'UP' THEN rs_dest.is_stop   -- swapped
					WHEN c.direction = 'DOWN' THEN rs_src.is_stop  -- swapped
			  	END DESC,
			rs_src.stop_sequence;`

		rows, err := pool.Query(context.Background(), query, temp_src+"%", temp_dest+"%")
		if err != nil {
			fmt.Println("error while finding the inactive route which absent in current_bus_route - ", err)
		}

		defer rows.Close()

		for rows.Next() {
			var route models.CurrentRoute
			rows.Scan(&route.BusId,
				&route.DriverId,
				&route.RouteId,
				&route.Direction,
				&route.RouteName,
				&route.Dest,
				&route.Src,
				&route.IsStop)

			FindStops(ctx, &route)
			fmt.Println("inactive route - ", route)
			Matched_routes = append(Matched_routes, route)
		}

	}
	return Matched_routes, nil
}

func ChangeRouteDirection(ctx context.Context, direction string) (bool, error) {
	query := `
				UPDATE current_bus_route AS cbr
				SET 
					direction = ar.direction,
					route_name = ar.route_name,
					src = ar.src,
					dest = ar.dest,
					route_id = ar.route_id
				FROM all_routes AS ar
				WHERE 
					ar.route_id = cbr.route_id
					AND ar.direction = $1
					AND cbr.direction <> $1;
			`
	_, err := pool.Exec(context.Background(), query, direction)

	if err != nil {
		fmt.Println("error while changing the routes - ", err)
		return false, err
	}
	return true, nil
}
