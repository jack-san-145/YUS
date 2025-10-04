package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	up_route_id, err1 := insert_route_to_db(up_route)
	down_route.Id = up_route_id //assign the up_route_id to the down_route_id
	_, err2 := insert_route_to_db(down_route)

	if err1 != nil && err2 != nil {
		return map[string]string{"status": "failed"}
	} else {
		return map[string]string{"status": "success"}
	}

}

// find the latest route id
func find_route_id() (int, error) {

	var route_id int
	query := "select route_id from all_routes where direction = 'UP' order by route_id desc limit 1"
	err := pool.QueryRow(context.Background(), query).Scan(&route_id)
	if errors.Is(err, sql.ErrNoRows) {
		return 1, nil
	} else if err != nil {
		fmt.Println("this error column is working")
		return -1, err
	}

	route_id += 1
	return route_id, nil
}

func insert_route_to_db(route *models.Route) (int, error) {

	var (
		arrival_time   string
		departure_time string
		route_name     string
		err            error
	)

	if route.Direction == "UP" {
		route.Id, err = find_route_id()
		if err != nil {
			fmt.Println("error while finding the route_id - ", err)
			return -1, err
		}

		route_name = route.UpRouteName
		arrival_time = route.ArrivalTime
		departure_time = route.UpDepartureTime
	} else if route.Direction == "DOWN" {
		route_name = route.DownRouteName
		arrival_time = route.ArrivalTime
		departure_time = route.DownDepartureTime
	}

	//inserting route and route_stops
	query := "insert into all_routes(route_id,route_name,src,dest,direction,departure_time,arrival_time) values($1,$2,$3,$4,$5,$6,$7);"
	_, err = pool.Exec(context.Background(), query, route.Id, route_name, route.Src, route.Dest, route.Direction, departure_time, arrival_time)
	if err != nil {
		fmt.Println("error while inserting route to db - ", err)
		return -1, fmt.Errorf("error")
	}

	for _, stop := range route.Stops {
		query = "insert into route_stops(route_id,route_name,direction,stop_name,is_stop,lat,lon,arrival_time,departure_time) values($1,$2,$3,$4,$5,$6,$7,$8)"
		_, err := pool.Exec(context.Background(), query, route.Id, route_name, route.Direction, stop.LocationName, stop.IsStop, stop.Lat, stop.Lon, stop.Arrival_time, stop.Departure_time)
		if err != nil {
			fmt.Println("error while inserting the route stops  - ", err)
			return -1, fmt.Errorf("error")
		}

	}
	return route.Id, nil
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
