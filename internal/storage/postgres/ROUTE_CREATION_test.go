package postgres

import (
	"testing"
	"time"
	"yus/internal/models"
)

func Test_SaveRoute_to_DB(t *testing.T) {

	Connect()
	route := models.Route{
		UpRouteName:       "thoothukudi to kcet",
		DownRouteName:     "kcet to thoothukudi",
		Src:               "kcet",
		Dest:              "thoothukudi",
		DownDepartureTime: "16:40",
		Stops: []models.RouteStops{
			{StopSequence: 1, LocationName: "sattur police station", Lat: "9.3540035", Lon: "77.9231079", IsStop: true, Arrival_time: "08:00"},
			{StopSequence: 2, LocationName: "sattur bus stand", Lat: "9.3538361", Lon: "77.9231022", IsStop: true, Arrival_time: "08:15"},
			{StopSequence: 3, LocationName: "sattur toll ghate", Lat: "9.3538361", Lon: "77.9231022", IsStop: false, Arrival_time: "08:35"},
			{StopSequence: 4, LocationName: "soolakrai", Lat: "9.3538361", Lon: "77.9231022", IsStop: true, Arrival_time: "08:35"},
			{StopSequence: 5, LocationName: "Kcet", Lat: "9.3538282", Lon: "77.9231091", IsStop: true, Arrival_time: "09:00"},
		},
		Created_At: time.Now(),
	}

	// Call your function
	status := SaveRoute_to_DB(&route)

	// ✅ Validate result
	t.Log("status - ", status)
}
