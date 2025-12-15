package postgres

import (
	"context"
	"log"
	"yus/internal/models"
)

func ScheduleBus(schedule *models.CurrentSchedule) error {

	query := `update current_bus_route set driver_id=1000 where driver_id=$1;
			  update current_bus_route set route_id=0 where route_id=$2;
			  update current_bus_route set driver_id=$3 where bus_id=$4;
			  update current_bus_route set route_id=$5 where bus_id=$6;
			`

	_, err := pool.Exec(context.Background(), query,
		schedule.DriverId,
		schedule.RouteId,
		schedule.DriverId,
		schedule.BusId,
		schedule.RouteId,
		schedule.BusId,
	)

	if err != nil {
		log.Println("error while schedule bus - ", err)
		return err
	}

	if schedule.RouteId != 0 {

		var route models.BusRoute

		route.BusID = schedule.BusId
		route.RouteId = schedule.RouteId
		route.Src, route.Dest, route.RouteName = find_src_dest_name_by_routeId(schedule.RouteId)

		go cache_this_route(&route)
	}
	return nil

}
