package postgres

import (
	"context"
	"fmt"
	"yus/internal/models"
)

func Store_new_driver_to_DB(new_driver *models.Driver, status *models.DriverAddedStatus) {
	query := "insert into drivers(name,mobile_no,email) values($1,$2,$3)"
	_, err := pool.Exec(context.Background(), query, new_driver.Name, new_driver.Mobile_no, new_driver.Email)
	if err != nil {
		fmt.Println("error while inserting the new driver - ", err)
		status.IsAdded = false
	}
	status.IsAdded = true
}
