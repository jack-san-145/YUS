package postgres

import (
	"context"
	"log"
	"yus/internal/models"

	"github.com/jackc/pgx/v5"
)

func ScheduleBus(schedule *models.CurrentSchedule) error {

	ctx := context.Background()
	tx, err := pool.Begin(ctx) //transaction for atomic operation
	if err != nil {
		log.Println("failed to begin tx:", err)
		return err
	}
	defer tx.Rollback(ctx) // safe rollback if commit doesn't happen

	batch := &pgx.Batch{}

	batch.Queue("update current_bus_route set driver_id=1000 where driver_id=$1;", schedule.DriverId)
	batch.Queue("update current_bus_route set route_id=0 where route_id=$1;", schedule.RouteId)
	batch.Queue("update current_bus_route set driver_id=$1,route_id=$2 where bus_id=$3;", schedule.DriverId, schedule.RouteId, schedule.BusId)

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Println("error while scheduling bus routes - ", err)
			return err
		}
	}
	err = tx.Commit(ctx) //close transaction
	if err != nil {
		log.Println("commit failed:", err)
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
