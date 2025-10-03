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

func Available_drivers() {
	query := "select driver_id,name from drivers"
	all_drivers, err := pool.Query(context.Background(), query)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("driver is empty - ", err)
		return
	} else if err != nil {
		fmt.Println("error while selecting the driver_id and name - ", err)
		return
	}

	defer all_drivers.Close()
	for all_drivers.Next() {
		// all_drivers.Scan()
	}

}
