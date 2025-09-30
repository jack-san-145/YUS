package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"yus/internal/models"
	"yus/internal/storage/postgres"
)

func Add_new_driver_handler(w http.ResponseWriter, r *http.Request) {
	var (
		Driver_array []models.Driver
		Status_array []models.DriverAddedStatus
	)
	err := json.NewDecoder(r.Body).Decode(&Driver_array)
	if err != nil {
		fmt.Println("error while decode the driver array - ", err)
		return
	}
	for _, driver := range Driver_array {
		var status models.DriverAddedStatus
		status.Name = driver.Name
		status.MobileNo = driver.Mobile_no
		status.Email = driver.Email

		fmt.Println("driver - ", driver)
		if validateMobileNo(driver.Mobile_no) && validateName(driver.Name) && validateEmail(driver.Email) {
			if postgres.Store_new_driver_to_DB(&driver) {
				status.IsAdded = true
			} else {
				status.IsAdded = false
			}

		} else {
			status.IsAdded = false
		}
		fmt.Println("status - ", status)
		Status_array = append(Status_array, status)
	}
	WriteJSON(w, r, Status_array)

}

// validate the name
func validateName(name string) bool {
	// Allows alphabets and spaces, 2 to 50 chars long
	re := regexp.MustCompile(`^[A-Za-z ]{2,50}$`)
	is_valid := re.MatchString(name)
	return is_valid
}

// validate the mobile_no with the regexp

func validateMobileNo(mobileNo string) bool {
	re := regexp.MustCompile(`^[6-9]\d{9}$`)
	is_valid := re.MatchString(mobileNo)
	return is_valid
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	is_valid := re.MatchString(email)
	return is_valid
}

// frontend -> backend: (post)

// {
// 	"email":"23ucs145@kamarajengg.edu.in",
// 	"password":"jack@145"
// }

// backend -> frontend
// {
// 	"login_status":"valid" or "Invalid"
// }
