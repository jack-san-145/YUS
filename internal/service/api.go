package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yus/internal/models"
)

func Reverse_Geocoding(current_location models.Location) {
	var (
		new_request *http.Request
		response    *http.Response
		result      models.NominatimResponse
		err         error
	)
	api := fmt.Sprintf(`https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json`, current_location.Latitude, current_location.Longitude)
	// response, err := http.Get(api)
	// if err != nil {

	// }
	// json.NewDecoder(response.Body).Decode()
	new_request, err = http.NewRequest("GET", api, nil)
	if err != nil {
		fmt.Println("error while creating the new request - ", err)
		return
	}
	new_request.Header.Set("User-Agent", "YUSreverseGeocoding/1.0 (yellohbus@gmail.com)")
	client := &http.Client{}
	response, err = client.Do(new_request)
	if err != nil {
		fmt.Println("error while get the reverse geocoding - ", err)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("error while decoding the response - ", err)
		return
	}

	fmt.Println("current location name - ", result.DisplayName)
	if result.Address.Village != "" {
		fmt.Println("village - ", result.Address.Village)
	} else if result.Address.Town != "" {
		fmt.Println("town - ", result.Address.Town)
	} else if result.Address.City != "" {
		fmt.Println("city - ", result.Address.City)
	}

}
