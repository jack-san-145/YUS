package postgres

import (
	"slices"
	"time"
	"yus/internal/models"
)

func SaveRoute_to_DB(up_route *models.Route) {
	up_route.Direction = "up"
	// down_route := find_down_route(up_route)

}

// func find_down_route(new_route *models.Route) *models.Route {

// }

func find_down_route(down_route models.Route) {
	down_route.Stops, down_route.DepartureTime, down_route.ArrivalTime = calculate_down_routeStops(&down_route)

}

// convert the string format to go's time.Time format to perform calcalution
func string_to_Gotime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}

// convert the go's time.Time foramt to string to store it to DB
func goTime_to_string(t time.Time) string {
	return t.Format("15:04")
}

// Calculate Down trip using segment durations from Up trip and 1-min stops
func calculate_down_routeStops(upRoute *models.Route) ([]models.RouteStops, string, string) {

	var (
		down_route_departure string
		down_route_arrival   string
	)

	upStops := upRoute.Stops
	//upStops = Madurai → Thirumangalam → Kallikudi → Kamaraj College

	downStops := upRoute.Stops //here just copy the upstops to downstops
	//upStops = Madurai → Thirumangalam → Kallikudi → Kamaraj College

	slices.Reverse(downStops) // to reverse only the upstops to downstops in-place

	// 1. Calculate segment durations from Up stop
	segmentDurations := []time.Duration{}
	for i := 0; i < len(upStops)-1; i++ {
		dep := string_to_Gotime(upStops[i].Departure_time)
		arr := string_to_Gotime(upStops[i+1].Arrival_time)

		//seg_duration = arrival_time - departure_time
		segmentDurations = append(segmentDurations, arr.Sub(dep))
	}

	/* | Segment                     | Duration |
	   | --------------------------- | -------- |
	   | Madurai → Thirumangalam     | 18 min   |
	   | Thirumangalam → Kallikudi   | 14 min   |
	   | Kallikudi → Kamaraj College | 13 min   |

	*/

	// therefore, segmentDurations = [18m, 14m, 13m]

	// 2. Initialize first stop's arr_time and depart_time for Down trip
	downDeparture := string_to_Gotime(upRoute.DownDepartureTime)
	downStops[0].Departure_time = goTime_to_string(downDeparture)
	downStops[0].Arrival_time = goTime_to_string(downDeparture)

	/*
		| Stop            | Arrival_time | Departure_time |
		| --------------- | ------------ | -------------- |
		| Kamaraj College | 16:40        | 16:40          |
		| Kallikudi       | ?            | ?              |
		| Thirumangalam   | ?            | ?              |
		| Madurai         | ?            | ?              |

	*/

	down_route_departure = goTime_to_string(downDeparture)
	currentTime := downDeparture

	// 3. reverse the segmentDurations slice
	slices.Reverse(segmentDurations)
	/*
		Before reverse: [18m, 14m, 13m]
		After reverse: [13m, 14m, 18m]
	*/

	// 4. Calculate arrival/departure for each stop
	for i := 0; i < len(downStops)-1; i++ {
		duration := segmentDurations[i]
		currentTime = currentTime.Add(duration)
		downStops[i+1].Arrival_time = goTime_to_string(currentTime)

		// Add 1-minute halt if stop
		if downStops[i+1].IsStop {
			currentTime = currentTime.Add(1 * time.Minute)
		}

		downStops[i+1].Departure_time = goTime_to_string(currentTime)
	}

	// assign the last stop Madurai's arrival time as final route departure time
	down_route_arrival = downStops[len(downStops)-1].Arrival_time

	return downStops, down_route_departure, down_route_arrival

	/*
	   1.Kamaraj College → Kallikudi

	   	Duration = 13 min
	   	Arrival = 16:40 + 13 min = 16:53
	   	Kallikudi IsStop = false → Departure = 16:53

	   2.Kallikudi → Thirumangalam

	   	Duration = 14 min
	   	Arrival = 16:53 + 14 min = 17:07
	   	Thirumangalam IsStop = true → +1 min halt → Departure = 17:08

	   3.Thirumangalam → Madurai

	   	Duration = 18 min
	   	Arrival = 17:08 + 18 min = 17:26
	   	Madurai IsStop = true → +1 min halt → Departure = 17:27

	   | Stop            | Arrival_time | Departure_time |
	   | --------------- | ------------ | -------------- |
	   | Kamaraj College | 16:40        | 16:40          |
	   | Kallikudi       | 16:53        | 16:53          |
	   | Thirumangalam   | 17:07        | 17:08          |
	   | Madurai         | 17:26        | 17:27          |
	*/
}
