package postgres

import (
	"context"
	"fmt"
	"yus/internal/models"
)

func (pg *PgStore) AssignDriverToBus(ctx context.Context, driverID int, busID int) error {
	// finds , is the driver_id existing in the current_bus_route
	is_driver_present, err := pg.DriverExistsInCBR(ctx, driverID)
	if err != nil {
		return err
	}

	if is_driver_present {

		// update the bus_route when the driver_id is exists
		err := pg.UpdateBusDriver(ctx, driverID, busID)
		if err != nil {
			return err
		}

	} else {
		// update the bus_route when the driver_id is not exists
		query := "update current_bus_route set driver_id = $1 where bus_id = $2"
		_, err = pool.Exec(context.Background(), query, driverID, busID)
		if err != nil {
			fmt.Println("error while update the driver_id for given bus_id - ", err)
			return err
		}
	}
	return nil
}

func (pg *PgStore) AssignRouteToBus(ctx context.Context, routeID int, busID int) error {

	var route models.BusRoute
	// finds , is the route_id existing in the current_bus_route
	is_route_present, err := pg.RouteExistsInCBR(ctx, routeID)
	if err != nil {
		return err
	}

	src, dest, route_name, _ := pg.GetSrcDestNameByRouteID(ctx, routeID)

	if is_route_present && routeID != 0 {

		route.RouteId, route.RouteName, route.Src, route.Dest, route.BusID = routeID, route_name, src, dest, busID
		// update the bus_route when the route_id is exists
		err := pg.UpdateBusRoute(ctx, &route)
		if err != nil {
			return err
		}

	} else {

		// update the bus_route when the route_id is not exists
		query := "update current_bus_route set route_id = $1 ,route_name = $2 ,src = $3 ,dest = $4 ,direction = 'UP' where bus_id = $5"
		_, err = pool.Exec(context.Background(), query, routeID, route_name, src, dest, busID)
		if err != nil {
			fmt.Println("error while update the route_id for given bus_id - ", err)
			return err
		}
	}
	if routeID != 0 {
		go pg.CacheRoute(ctx, &route)
	}
	return nil
}

func (pg *PgStore) GetSrcDestNameByRouteID(ctx context.Context, routeID int) (string, string, string, error) {
	var (
		route_name string
		src        string
		dest       string
	)
	query := "select route_name,src,dest from all_routes where route_id = $1 and direction = 'UP' "
	err := pool.QueryRow(context.Background(), query, routeID).Scan(&route_name, &src, &dest)
	if err != nil {
		fmt.Println("error while finding the route_name,src and dest - ", err)
		return "", "", "", err
	}
	return src, dest, route_name, nil
}

// finds , is the route_id existing in the current_bus_route
func (pg *PgStore) RouteExistsInCBR(ctx context.Context, routeID int) (bool, error) {
	var is_route_present bool
	query := "select exists(select 1 from current_bus_route where route_id = $1)"
	err := pool.QueryRow(context.Background(), query, routeID).Scan(&is_route_present)
	if err != nil {
		fmt.Println("error while finding the existance of route_id in current_bus_route - ", err)
		return false, err
	}
	return is_route_present, nil
}

// finds , is the route_id existing in the current_bus_route
func (pg *PgStore) DriverExistsInCBR(ctx context.Context, driverID int) (bool, error) {
	var is_driver_present bool
	query := "select exists(select 1 from current_bus_route where driver_id = $1)"
	err := pool.QueryRow(context.Background(), query, driverID).Scan(&is_driver_present)
	if err != nil {
		fmt.Println("error while finding the existance of driver_id in current_bus_route - ", err)
		return false, err
	}
	return is_driver_present, nil
}

// update the bus_route when the route_id is exists on current_bus_route
func (pg *PgStore) UpdateBusRoute(ctx context.Context, route *models.BusRoute) error {

	//set the already mapped other bus's route_id as 0
	query := "update current_bus_route set route_id = 0,route_name ='',src ='',dest ='' where route_id = $1"
	_, err := pool.Exec(context.Background(), query, route.RouteId)
	if err != nil {
		fmt.Println("error while update the existing route as 0 - ", err)
		return err
	}

	//map the new route with bus
	query = "update current_bus_route set route_id = $1 ,route_name = $2 ,src = $3 ,dest = $4 where bus_id = $5"
	_, err = pool.Exec(context.Background(), query, route.RouteId, route.RouteName, route.Src, route.Dest, route.BusID)
	if err != nil {
		fmt.Println("error while update the route_id for given bus_id - ", err)
		return err
	}
	return nil
}

// update the bus_route when the route_id is exists on current_bus_route
func (pg *PgStore) UpdateBusDriver(ctx context.Context, driverID int, busID int) error {

	//set the already mapped other bus's route_id as 0
	query := "update current_bus_route set driver_id = 1000 where driver_id = $1"
	_, err := pool.Exec(context.Background(), query, driverID)
	if err != nil {
		fmt.Println("error while update the existing driver as 1000 - ", err)
		return err
	}

	//map the new driver with bus
	query = "update current_bus_route set driver_id = $1 where bus_id = $2"
	_, err = pool.Exec(context.Background(), query, driverID, busID)
	if err != nil {
		fmt.Println("error while update the driver_id for given bus_id - ", err)
		return err
	}
	return nil
}

func (pg *PgStore) AddBus(ctx context.Context, busID int) error {

	var is_bus_exists bool
	fmt.Println("bus_id - ", busID)

	//chech if the bus already exists or not in current_bus_route
	query := "select exists(select 1 from current_bus_route where bus_id = $1)"
	err := pool.QueryRow(context.Background(), query, busID).Scan(&is_bus_exists)
	if err != nil {
		fmt.Println("error while finding the existance of the bus - ", err)
		return fmt.Errorf("failed")
	}

	//if exists returns "bus already exists"
	if is_bus_exists {
		return fmt.Errorf("bus already exists")
	}

	//if doesn't exists then add the new bus to current_bus_route
	query = "insert into current_bus_route(bus_id) values($1) "
	_, err = pool.Exec(context.Background(), query, busID)
	if err != nil {
		fmt.Println("error while adding new bus - ", err)
		return fmt.Errorf("failed")
	}
	return nil
}

func (pg *PgStore) CacheRoute(ctx context.Context, route *models.BusRoute) error {
	var is_route_exists bool

	//check if the route already exist in cached_bus_route or not
	query := "select exists(select 1 from cached_bus_route where bus_id = $1 and route_id = $2)"
	err := pool.QueryRow(context.Background(), query, route.BusID, route.RouteId).Scan(&is_route_exists)
	if err != nil {
		fmt.Println("error while finding the existance of the bus_route in cached_bus_route - ", err)
		return err
	}

	// if it doesn't exists then add this bus_route to cached_bus_route otherwise do nothing
	if !is_route_exists {
		query = "insert into cached_bus_route(bus_id,route_id,route_name,src,dest) values($1,$2,$3,$4,$5)"
		_, err = pool.Exec(context.Background(), query, route.BusID, route.RouteId, route.RouteName, route.Src, route.Dest)
		if err != nil {
			fmt.Println("error while insert the route to the cached route - ", err)
			return err
		}
	}
	return nil
}
