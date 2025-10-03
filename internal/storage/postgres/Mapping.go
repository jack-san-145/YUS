package postgres

import (
	"context"
	"fmt"
)

func Map_driver_with_bus(driver_id int, bus_id int) {

}

func Map_route_with_bus(route_id int, bus_id int) error {

	// finds , is the route_id existing in the current_bus_route
	is_route_present, err := find_route_exists_CBR(route_id)
	if err != nil {
		return err
	}

	if is_route_present {

		// update the bus_route when the route_id is exists
		err := update_bus_route(route_id, bus_id)
		if err != nil {
			return err
		}

	} else {
		// update the bus_route when the route_id is not exists
		query := "update current_bus_route set route_id = $1 where bus_id = $2"
		_, err = pool.Exec(context.Background(), query, route_id, bus_id)
		if err != nil {
			fmt.Println("error while update the route_id for given bus_id - ", err)
			return err
		}
	}
	go cache_this_route(route_id, bus_id)
	return nil
}

// finds , is the route_id existing in the current_bus_route
func find_route_exists_CBR(route_id int) (bool, error) {
	var is_route_present bool
	query := "select exists(select 1 from current_bus_route where route_id = $1)"
	err := pool.QueryRow(context.Background(), query, route_id).Scan(&is_route_present)
	if err != nil {
		fmt.Println("error while finding the existance of route_id in current_bus_route - ", err)
		return false, err
	}
	return is_route_present, nil
}

// update the bus_route when the route_id is exists on current_bus_route
func update_bus_route(route_id int, bus_id int) error {

	//set the already mapped other bus's route_id as 0
	query := "update current_bus_route set route_id = 0 where route_id = $1"
	_, err := pool.Exec(context.Background(), query, route_id)
	if err != nil {
		fmt.Println("error while update the existing route as 0 - ", err)
		return err
	}

	//map the new route with bus
	query = "update current_bus_route set route_id = $1 where bus_id = $2"
	_, err = pool.Exec(context.Background(), query, route_id, bus_id)
	if err != nil {
		fmt.Println("error while update the route_id for given bus_id - ", err)
		return err
	}
	return nil
}

func Add_new_bus(bus_id int) error {

	var is_bus_exists bool
	fmt.Println("bus_id - ", bus_id)

	//chech if the bus already exists or not in current_bus_route
	query := "select exists(select 1 from current_bus_route where bus_id = $1)"
	err := pool.QueryRow(context.Background(), query, bus_id).Scan(&is_bus_exists)
	if err != nil {
		fmt.Println("error while finding the existance of the bus - ", err)
		return fmt.Errorf("failed")
	}

	//if exists returns "bus already exists"
	if is_bus_exists {
		return fmt.Errorf("bus already exists")
	}

	//if doesn't exists then add the new bus to current_bus_route
	query = "insert into current_bus_route(route_id,bus_id) values(0,$1) "
	_, err = pool.Exec(context.Background(), query, bus_id)
	if err != nil {
		fmt.Println("error while adding new bus - ", err)
		return fmt.Errorf("failed")
	}
	return nil
}

func cache_this_route(route_id int, bus_id int) {
	var is_route_exists bool

	//check if the route already exist in cached_bus_route or not
	query := "select exists(select 1 from cached_bus_route where bus_id = $1 and route_id = $2)"
	err := pool.QueryRow(context.Background(), query, bus_id, route_id).Scan(&is_route_exists)
	if err != nil {
		fmt.Println("error while finding the existance of the bus_route in cached_bus_route - ", err)
	}

	// if it doesn't exists then add this bus_route to cached_bus_route otherwise do nothing
	if !is_route_exists {
		query = "insert into cached_bus_route(bus_id,route_id) values($1,$2)"
		_, err = pool.Exec(context.Background(), query, bus_id, route_id)
		if err != nil {
			fmt.Println("error while insert the route to the cached route - ", err)
		}
	}

}
