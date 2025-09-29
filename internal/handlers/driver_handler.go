package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
	"yus/internal/storage/postgres"
)

func Add_new_driver_handler(w http.ResponseWriter, r *http.Request) {
	var Driver_array []models.Driver
	err := json.NewDecoder(r.Body).Decode(&Driver_array)
	if err != nil {
		fmt.Println("error while decode the driver array - ", err)
		return
	}
	for _, driver := range Driver_array {
		fmt.Println("driver - ", driver)
		go postgres.Store_new_driver_to_DB(&driver)
	}

}

func validate_mobileno(mobile_no string) {
	re := ``
}
