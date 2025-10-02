package postgres

import (
	"context"
	"fmt"
)

func Map_route_with_bus(route_id int, bus_id int) error {

	var (
		is_bus_present bool
		// is_route_present bool
	)

	ctx := context.Background()

	is_route_present, err := find_route_exists_CBR(route_id)
	if err != nil {
		return err
	}

	if is_bus_present {

	}

	if is_route_present {
		query := "update current_bus_route set route_id = 0 where route_id = $1"
		_, err = pool.Exec(ctx, query, route_id)
		if err != nil {
			fmt.Println("error while update the existing route as 0 - ", err)
			return fmt.Errorf(err.Error())
		}

	}

	query := "select exists(select 1 from current_bus_route where bus_id = $1)"
	pool.QueryRow(ctx, query, bus_id).Scan(&is_bus_present)
	if is_bus_present {
		query = "update current_bus_route set route_id = "
	}
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

// func set_existing

func Add_new_bus(bus_no int) error {
	query := "insert into current_bus_route(bus_id) values($1) "
	_, err := pool.Exec(context.Background(), query, bus_no)
	if err != nil {
		fmt.Println("error while adding new bus - ", err)
		return fmt.Errorf("failed")
	}
	return nil
}
