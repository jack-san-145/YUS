package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func SaveRoute_to_DB(up_route *models.Route) map[string]string {

	up_route.Direction = "UP"
	services.Calculate_Uproute_departure(up_route)
	down_route := services.Find_down_route(*up_route)
	fmt.Println()
	fmt.Println("up route - ", up_route)
	fmt.Println()
	fmt.Println("down route - ", down_route)

	//check is the new  route exist or not ?

	err := check_if_route_exist(up_route.Src, up_route.Dest, up_route.Stops)
	if err != nil {
		fmt.Println(err.Error())
		return map[string]string{"status": err.Error()}
	}

	fmt.Println("going to insert route to table")
	//inserting both up and down routes to db
	err1 := insert_route_to_db(up_route)
	err2 := insert_route_to_db(down_route)

	if err1 != nil && err2 != nil {
		return map[string]string{"status": "failed"}
	} else {
		return map[string]string{"status": "success"}
	}

}

func insert_route_to_db(route *models.Route) error {

	var (
		arrival_time   string
		departure_time string
	)

	if route.Direction == "UP" {
		arrival_time = route.ArrivalTime
		departure_time = route.UpDepartureTime
	} else if route.Direction == "DOWN" {
		arrival_time = route.ArrivalTime
		departure_time = route.DownDepartureTime
	}

	//inserting route and route_stops
	query := "insert into all_routes(src,dest,direction,departure_time,arrival_time) values($1,$2,$3,$4,$5) returning route_id;"
	err := pool.QueryRow(context.Background(), query, route.Src, route.Dest, route.Direction, departure_time, arrival_time).Scan(&route.Id)
	if err != nil {
		fmt.Println("error while inserting route to db - ", err)
		return fmt.Errorf("error")
	}

	for _, stop := range route.Stops {
		query = "insert into route_stops(route_id,stop_name,is_stop,lat,lon,arrival_time,departure_time) values($1,$2,$3,$4,$5,$6,$7)"
		_, err := pool.Exec(context.Background(), query, route.Id, stop.LocationName, stop.IsStop, stop.Lat, stop.Lon, stop.Arrival_time, stop.Departure_time)
		if err != nil {
			fmt.Println("error while inserting the route stops  - ", err)
			return fmt.Errorf("error")
		}

	}
	return nil
}

// to check if the route is exist on DB or not
func check_if_route_exist(src string, dest string, stops []models.RouteStops) error {

	var is_match_found_inthis_routes bool

	query := "select route_id from all_routes where src = $1 and dest = $2 ;"
	route_id_rows, err_error := pool.Query(context.Background(), query, src, dest)
	if err_error != nil {
		fmt.Println("error while finding the route id - ", err_error)
	}

	for route_id_rows.Next() {
		var (
			route_id int
		)

		route_id_rows.Scan(&route_id)

		if route_id == 0 {
			return nil
		}
		query = "select stop_name,is_stop from route_stops where route_id = $1"
		rows, err := pool.Query(context.Background(), query, route_id)
		if err != nil {
			fmt.Println("error while accesing the stopname - ", err)
		}
		if err == sql.ErrNoRows {
			fmt.Println("no rows")
			continue
		}

		for _, stop := range stops {
			var (
				stop_name string
				is_stop   bool
			)
			if !rows.Next() {
				break
			}

			err := rows.Scan(&stop_name, &is_stop)
			if err != nil {
				fmt.Println("error while accessing the stopname and is_stop - ", err)
			}

			if !(stop_name == stop.LocationName && is_stop == stop.IsStop) {
				is_match_found_inthis_routes = false
				break
			} else {
				is_match_found_inthis_routes = true
			}

		}
		rows.Close()
		if is_match_found_inthis_routes {
			return fmt.Errorf("root already exists")
		}
	}
	route_id_rows.Close()
	return nil
}
