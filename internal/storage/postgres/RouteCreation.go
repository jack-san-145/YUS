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

	//inserting up route
	query := "insert into all_routes(src,dest,direction,departure_time,arrival_time) values($1,$2,$3,$4,$5);"
	_, err := pool.Exec(context.Background(), query, up_route.Src, up_route.Dest, up_route.Direction, up_route.UpDepartureTime, up_route.ArrivalTime)
	if err != nil {
		fmt.Println("error while inserting up route to the db - ", err)
		return map[string]string{"status": "failed"}
	}

	//inserting down route
	query = "insert into all_routes(src,dest,direction,departure_time,arrival_time) values($1,$2,$3,$4,$5);"
	_, err = pool.Exec(context.Background(), query, down_route.Src, down_route.Dest, down_route.Direction, down_route.DownDepartureTime, down_route.ArrivalTime)
	if err != nil {
		fmt.Println("error while inserting down route to the db - ", err)
		return map[string]string{"status": "failed"}
	}
	return map[string]string{"status": "success"}

}

func check_if_route_exist(src string, dest string, stops []models.RouteStops) error {
	var route_id int
	query := "select coalesce( (select route_id from all_routes where src = $1 and dest = $2),0) as route_id;"
	pool.QueryRow(context.Background(), query, src, dest).Scan(&route_id)

	if route_id == 0 {
		return nil
	}
	query = "select stop_name,is_stop from route_stops where route_id = $1"
	rows, err := pool.Query(context.Background(), query, route_id)
	if err == sql.ErrNoRows {
		fmt.Println("no rows")
		return nil
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
			return nil
		}
	}

	return fmt.Errorf("Root already exists")
}
