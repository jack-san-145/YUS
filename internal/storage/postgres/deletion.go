package postgres

import (
	"context"
	"fmt"
)

func RemoveRoute(ctx context.Context, routeID int) error {
	query := fmt.Sprintf(`update current_bus_route set route_id=0, direction='', route_name='', src='', dest='' where route_id = %d;
							delete from cached_bus_route where route_id = %d;
							delete from route_stops where route_id = %d;  
							delete from all_routes where route_id = %d; 
						`, routeID, routeID, routeID, routeID)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		fmt.Println("error while deleting route - ", err)
		return err
	}
	return nil
}

func Remove_bus(bus_id int) error {

	query := fmt.Sprintf(`delete from cached_bus_route where bus_id = %d;
							delete from current_bus_route where bus_id = %d;
						`, bus_id, bus_id)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while removing bus - ", err)
		return err
	}
	return nil
}

func Remove_driver(driver_id int) error {
	query := fmt.Sprintf(`update current_bus_route set driver_id = 1000 where driver_id = %d ;
							delete from drivers where driver_id = %d ;
						`, driver_id, driver_id)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while removing driver - ", err)
		return err
	}
	return nil

}
