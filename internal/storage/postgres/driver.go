package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
)

func Store_new_driver_to_DB(new_driver *models.Driver) bool {
	query := "insert into drivers(name,mobile_no,email) values($1,$2,$3)"
	_, err := pool.Exec(context.Background(), query, new_driver.Name, new_driver.Mobile_no, new_driver.Email)
	if err != nil {
		fmt.Println("error while inserting the new driver - ", err)
		return false
	}
	fmt.Println("working")
	return true
}

func Available_drivers() []models.AvailableDriver {
	var all_available_drivers []models.AvailableDriver
	query := "select driver_id,name from drivers"
	all_drivers, err := pool.Query(context.Background(), query)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("driver is empty - ", err)
		return nil
	} else if err != nil {
		fmt.Println("error while selecting the driver_id and name - ", err)
		return nil
	}

	defer all_drivers.Close()
	for all_drivers.Next() {
		var driver models.AvailableDriver
		err = all_drivers.Scan(&driver.Id, &driver.Name)
		if err != nil {
			fmt.Println("error while scanning the driver from DB - ", err)
			continue
		}
		query = "select exists (select 1 from current_bus_route where driver_id = $1 ) "
		err := pool.QueryRow(context.Background(), query, driver.Id).Scan(&driver.Available)
		if err != nil {
			fmt.Println("error while checking existance of the driver in current_bus_route - ", err)
			continue
		}
		all_available_drivers = append(all_available_drivers, driver)
	}
	return all_available_drivers
}
