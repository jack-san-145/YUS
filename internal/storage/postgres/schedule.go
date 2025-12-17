package postgres

import (
	"context"
	"log"
	"yus/internal/models"

	"github.com/jackc/pgx/v5"
)

func (pg *PgStore) ScheduleBus(ctx context.Context, schedule *models.CurrentSchedule) error {

	var route models.BusRoute

	route.BusID = schedule.BusId
	route.RouteId = schedule.RouteId
	route.Src, route.Dest, route.RouteName, _ = pg.GetSrcDestNameByRouteID(ctx, schedule.RouteId)

	tx, err := pg.Pool.Begin(ctx) //transaction for atomic operation
	if err != nil {
		log.Println("failed to begin tx:", err)
		return err
	}
	defer tx.Rollback(ctx) // safe rollback if commit doesn't happen

	batch := &pgx.Batch{}

	batch.Queue("update current_bus_route set driver_id=1000 where driver_id=$1;", schedule.DriverId)
	batch.Queue("update current_bus_route set route_id=0 where route_id=$1;", schedule.RouteId)
	batch.Queue("update current_bus_route set driver_id = $1,route_id = $2,route_name = $3,src = $4,dest = $5 where bus_id = $6;",
		schedule.DriverId,
		schedule.RouteId,
		route.RouteName,
		route.Src,
		route.Dest,
		schedule.BusId,
	)

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Println("error while scheduling bus routes - ", err)
			br.Close()
			return err
		}
	}

	br.Close()

	err = tx.Commit(ctx) //close transaction
	if err != nil {
		log.Println("commit failed:", err)
		return err
	}

	if schedule.RouteId != 0 {

		go pg.CacheRoute(context.Background(), &route)
	}
	return nil

}
