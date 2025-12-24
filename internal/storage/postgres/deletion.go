package postgres

import (
	"context"
	"fmt"
	"log"
	"yus/internal/models"
)

func (pg *PgStore) RemoveRoute(ctx context.Context, routeID int) error {
	query := fmt.Sprintf(`update current_bus_route set route_id=0, direction='', route_name='', src='', dest='' where route_id = %d;
							delete from cached_bus_route where route_id = %d;
							delete from route_stops where route_id = %d;  
							delete from all_routes where route_id = %d; 
						`, routeID, routeID, routeID, routeID)
	_, err := pg.Pool.Exec(ctx, query)
	if err != nil {
		fmt.Println("error while deleting route - ", err)
		return err
	}
	return nil
}

func (pg *PgStore) RemoveBus(ctx context.Context, busID int) error {

	query := fmt.Sprintf(`delete from cached_bus_route where bus_id = %d;
							delete from current_bus_route where bus_id = %d;
						`, busID, busID)
	_, err := pg.Pool.Exec(ctx, query)
	if err != nil {
		fmt.Println("error while removing bus - ", err)
		return err
	}
	return nil
}

func (pg *PgStore) RemoveDriver(ctx context.Context, driverID int, mode string) error {
	err := pg.ClearDriverRemovalRequest(ctx, driverID) //just remove the driver removal request from db
	if err != nil {
		return err
	}

	if mode != "OK" {
		return nil
	}

	query := fmt.Sprintf(`update current_bus_route set driver_id = 1000 where driver_id = %d ;
							delete from drivers where driver_id = %d ;
						`, driverID, driverID)
	_, err = pg.Pool.Exec(ctx, query)
	if err != nil {
		fmt.Println("error while removing driver - ", err)
		return err
	}
	return nil

}

func (pg *PgStore) StoreDriverRemovalRequest(ctx context.Context, driverID int) error {
	query := "insert into driver_removal_request(driver_id) values ($1);"
	_, err := pg.Pool.Exec(ctx, query, driverID)
	if err != nil {
		log.Println("error while store driver removal request - ", err)
		return err
	}
	return nil
}

func (pg *PgStore) GetDriverRemovalRequest(ctx context.Context) ([]models.DriverRemovalRequest, error) {
	var Allrequests []models.DriverRemovalRequest
	query := "select * from driver_removal_request"
	rows, err := pg.Pool.Query(ctx, query)
	if err != nil {
		log.Println("error while get driver removal requests - ", err)
		return Allrequests, err
	}
	defer rows.Close()

	for rows.Next() {
		var request models.DriverRemovalRequest
		err := rows.Scan(&request.DriverId,
			&request.Created_At)
		if err != nil {
			log.Println("error while scan the driver removal request - ", err)
			return Allrequests, err
		}
		Allrequests = append(Allrequests, request)
	}
	return Allrequests, nil
}

func (pg *PgStore) ClearDriverRemovalRequest(ctx context.Context, driverID int) error {
	query := "delete from driver_removal_request where driver_id = $1;"
	_, err := pg.Pool.Exec(ctx, query, driverID)
	if err != nil {
		log.Println("error while clear driver removal request from db - ", err)
		return err
	}
	return nil
}
