package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func Get_Allotted_Bus(driver_id int) (models.AllotedBus, error) {

	var alloted_bus models.AllotedBus
	query := "select bus_id,route_id,route_name,direction,src,dest from current_bus_route where driver_id = $1"
	err := pool.QueryRow(context.Background(), query, driver_id).Scan(&alloted_bus.BusID,

		&alloted_bus.RouteId,
		&alloted_bus.RouteName,
		&alloted_bus.Direction,
		&alloted_bus.Src,
		&alloted_bus.Dest)

	alloted_bus.RouteName = services.Convert_to_Normal(alloted_bus.RouteName)
	alloted_bus.Src = services.Convert_to_Normal(alloted_bus.Src)
	alloted_bus.Dest = services.Convert_to_Normal(alloted_bus.Dest)

	alloted_bus.DriverId = driver_id
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("no bus is allotted for driver_id - %v", driver_id)
		return alloted_bus, fmt.Errorf("no bus alloted")
	} else if err != nil {
		fmt.Println("error while finding the allotted bus for driver - ", err)
		return alloted_bus, err
	}
	return alloted_bus, nil
}

func Store_new_driver_to_DB(new_driver *models.Driver) error {
	query := "insert into drivers(driver_name,mobile_no) values($1,$2)"
	_, err := pool.Exec(context.Background(), query, new_driver.Name, new_driver.Mobile_no)
	if err != nil {
		fmt.Println("error while inserting the new driver - ", err)
		return err
	}
	fmt.Println("working")
	return nil
}

func Check_Driver_exits(driver_id int) (bool, error) {
	var exists bool
	query := "select exists(select 1 from drivers where driver_id = $1)"
	err := pool.QueryRow(context.Background(), query, driver_id).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking the existance of driver - ", err)
		return exists, err
	}
	return exists, nil
}

func Set_driver_password(driver_id int, driver_email string, password string) (bool, error) {
	hashed_pass := services.Hash_this_password(password)
	if exists, _ := Check_Driver_exits(driver_id); exists {
		query := "update drivers set password = $1,email = $2 where driver_id = $3 "
		_, err := pool.Exec(context.Background(), query, hashed_pass, driver_email, driver_id)
		if err != nil {
			fmt.Println("error while update the driver's password - ", err)
			return false, err
		}
	}
	return true, nil
}

func ValidateDriver(driver_id int, pass string) bool {

	var DB_pass string
	query := "select password from drivers where driver_id = $1"
	err := pool.QueryRow(context.Background(), query, driver_id).Scan(&DB_pass)
	if err != nil {
		fmt.Println("error while validate the driver - ", err)
		return false
	}

	if services.Is_password_matched(DB_pass, pass) {
		return true
	}
	return false
}

func Available_drivers() []models.AvailableDriver {
	var (
		all_available_drivers []models.AvailableDriver
		is_driver_exists      bool
	)

	query := "select driver_id,driver_name,mobile_no from drivers"
	all_drivers, err := pool.Query(context.Background(), query)
	if err != nil {
		fmt.Println("error while selecting the driver_id , driver_name and mobile no - ", err)
		return nil
	}

	defer all_drivers.Close()
	for all_drivers.Next() {
		var driver models.AvailableDriver
		err = all_drivers.Scan(&driver.Id, &driver.Name, &driver.MobileNo)
		if err != nil {
			fmt.Println("error while scanning the driver from DB - ", err)
			continue
		}
		query = "select exists (select 1 from current_bus_route where driver_id = $1 ) "
		err := pool.QueryRow(context.Background(), query, driver.Id).Scan(&is_driver_exists)
		if err != nil {
			fmt.Println("error while checking existance of the driver in current_bus_route - ", err)
			continue
		}
		if !is_driver_exists {
			driver.Available = true
		}
		all_available_drivers = append(all_available_drivers, driver)
	}
	return all_available_drivers
}
