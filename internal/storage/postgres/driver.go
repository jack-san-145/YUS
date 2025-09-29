package postgres

import (
	"context"
	"fmt"
	"yus/internal/models"
)

func Store_new_driver_to_DB(new_driver *models.Driver) {
	query := "insert into drivers(name,mobile_no) values($1,$2)"
	_, err := pool.Exec(context.Background(), query, new_driver.Name, new_driver.Mobile_no)
	if err != nil {
		fmt.Println("error while inserting the new driver - ", err)
		return
	}

}
