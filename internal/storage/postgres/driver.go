package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func (pg *PgStore) GetAllottedBusForDriver(ctx context.Context, driverID int) (models.AllotedBus, error) {

	var alloted_bus models.AllotedBus
	query := "select bus_id,route_id,route_name,direction,src,dest from current_bus_route where driver_id = $1"
	err := pg.Pool.QueryRow(ctx, query, driverID).Scan(&alloted_bus.BusID,

		&alloted_bus.RouteId,
		&alloted_bus.RouteName,
		&alloted_bus.Direction,
		&alloted_bus.Src,
		&alloted_bus.Dest)

	alloted_bus.RouteName = services.Convert_to_Normal(alloted_bus.RouteName)
	alloted_bus.Src = services.Convert_to_Normal(alloted_bus.Src)
	alloted_bus.Dest = services.Convert_to_Normal(alloted_bus.Dest)

	alloted_bus.DriverId = driverID
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("no bus is allotted for driver_id - %v", driverID)
		return alloted_bus, fmt.Errorf("no bus alloted")
	} else if err != nil {
		fmt.Println("error while finding the allotted bus for driver - ", err)
		return alloted_bus, err
	}
	return alloted_bus, nil
}

func (pg *PgStore) AddDriver(ctx context.Context, driver *models.Driver) error {
	query := "insert into drivers(driver_name,mobile_no) values($1,$2)"
	_, err := pg.Pool.Exec(ctx, query, driver.Name, driver.Mobile_no)
	if err != nil {
		fmt.Println("error while inserting the new driver - ", err)
		return err
	}
	fmt.Println("working")
	return nil
}

func (pg *PgStore) DriverExists(ctx context.Context, driverID int) (bool, error) {
	var exists bool
	query := "select exists(select 1 from drivers where driver_id = $1)"
	err := pg.Pool.QueryRow(ctx, query, driverID).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking the existance of driver - ", err)
		return exists, err
	}
	return exists, nil
}

func (pg *PgStore) SetDriverPassword(ctx context.Context, driverID int, email string, password string) (bool, error) {
	hashed_pass := services.Hash_this_password(password)
	if exists, _ := pg.DriverExists(ctx, driverID); exists {
		query := "update drivers set password = $1,email = $2 where driver_id = $3 "
		_, err := pg.Pool.Exec(ctx, query, hashed_pass, email, password)
		if err != nil {
			fmt.Println("error while update the driver's password - ", err)
			return false, err
		}
	}
	return true, nil
}

func (pg *PgStore) ValidateDriver(ctx context.Context, driverID int, password string) (bool, error) {

	var DB_pass string
	query := "select password from drivers where driver_id = $1"
	err := pg.Pool.QueryRow(ctx, query, driverID).Scan(&DB_pass)
	if err != nil {
		fmt.Println("error while validate the driver - ", err)
		return false, err
	}

	if services.Is_password_matched(DB_pass, password) {
		return true, nil
	}
	return false, nil
}

func (pg *PgStore) GetAvailableDrivers(ctx context.Context) ([]models.AvailableDriver, error) {
	var (
		all_available_drivers []models.AvailableDriver
		is_driver_exists      bool
	)

	query := "select driver_id,driver_name,mobile_no from drivers"
	all_drivers, err := pg.Pool.Query(ctx, query)
	if err != nil {
		fmt.Println("error while selecting the driver_id , driver_name and mobile no - ", err)
		return nil, err
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
		err := pg.Pool.QueryRow(ctx, query, driver.Id).Scan(&is_driver_exists)
		if err != nil {
			fmt.Println("error while checking existance of the driver in current_bus_route - ", err)
			continue
		}
		if !is_driver_exists {
			driver.Available = true
		}
		all_available_drivers = append(all_available_drivers, driver)
	}
	return all_available_drivers, nil
}
