package services

import (
	"slices"
	"time"
	"yus/internal/models"
)

// Admin provides only arrival_time for each stop
// here I calculates departure_time for each stop -> calculated by arrival_time
func Calculate_Uproute_departure(upRoute *models.Route) {
	upStops := upRoute.Stops

	//convert the normal names to camel-case to store in DB
	upRoute.UpRouteName = Convert_to_CamelCase(upRoute.UpRouteName)
	upRoute.Src = Convert_to_CamelCase(upRoute.Src)
	upRoute.Dest = Convert_to_CamelCase(upRoute.Dest)

	for i := 0; i < len(upStops); i++ {
		if i == 0 {
			// First stop: departure = arrival
			upStops[i].Departure_time = upStops[i].Arrival_time
		} else if i == len(upStops)-1 {
			//last stop: departure = arrival , it don't want departure bcz its the destination
			upStops[i].Departure_time = upStops[i].Arrival_time
		} else {
			// Other stops: departure = arrival + 1 min if IsStop
			arrival := string_to_Gotime(upStops[i].Arrival_time)
			if upStops[i].IsStop {
				upStops[i].Departure_time = goTime_to_string(arrival.Add(1 * time.Minute))
			} else {
				upStops[i].Departure_time = upStops[i].Arrival_time
			}
		}
		//convert the stop names(normal names )to camel-case to store in DB
		upStops[i].LocationName = Convert_to_CamelCase(upStops[i].LocationName)
	}
	upRoute.Stops = upStops
	upRoute.UpDepartureTime = upStops[0].Departure_time
	upRoute.ArrivalTime = upStops[len(upStops)-1].Arrival_time // reaches the destination

}

func Find_down_route(down_route models.Route) *models.Route {
	calculate_down_routeStops(&down_route)
	down_route.Src, down_route.Dest = down_route.Dest, down_route.Src
	down_route.Direction = "DOWN"
	return &down_route

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

// Down route is fully derived from Up route segment durations
// Admin provides only DownDepartureTime (first stop of Down route) <=> calculates by departure_time
func calculate_down_routeStops(down_route *models.Route) {

	// Step 1: Initialize Down route stops
	downStops := make([]models.RouteStops, len(down_route.Stops))
	copy(downStops, down_route.Stops)

	//convert the normal names to camel-case to store in DB
	down_route.DownRouteName = Convert_to_CamelCase(down_route.DownRouteName)
	down_route.Src = Convert_to_CamelCase(down_route.Src)
	down_route.Dest = Convert_to_CamelCase(down_route.Dest)

	// Step 2: Calculate segment durations from Up route
	segmentDurations := []time.Duration{}
	for i := 0; i < len(downStops)-1; i++ {
		dep := string_to_Gotime(downStops[i].Departure_time)
		arr := string_to_Gotime(downStops[i+1].Arrival_time)
		segmentDurations = append(segmentDurations, arr.Sub(dep))
	}

	/*
	   Example Up segments (Up route):
	   | Segment                     | Duration |
	   | --------------------------- | -------- |
	   | Madurai → Thirumangalam     | 18 min   |
	   | Thirumangalam → Kallikudi   | 14 min   |
	   | Kallikudi → Kamaraj College | 13 min   |
	   segmentDurations = [18m, 14m, 13m]
	*/

	//step:3

	slices.Reverse(downStops) // reverse Up stops → Down stops

	// Step 4: Set first stop arrival & departure
	downDeparture := string_to_Gotime(down_route.DownDepartureTime)
	downStops[0].Arrival_time = down_route.DownDepartureTime
	downStops[0].Departure_time = down_route.DownDepartureTime
	currentTime := downDeparture
	/*
	   | Stop            | Arrival | Departure |
	   | --------------- | ------- | --------- |
	   | Kamaraj College | 16:40   | 16:40     |
	*/

	// Step 5: Reverse segment durations for Down trip
	slices.Reverse(segmentDurations)

	/*
	   segmentDurations before: [18m, 14m, 13m]
	   segmentDurations after reverse: [13m, 14m, 18m]
	*/

	// Step 6: Calculate arrival/departure for each stop
	for i := 0; i < len(downStops)-1; i++ {
		duration := segmentDurations[i]
		// Arrival at next stop = previous departure + segment duration
		currentTime = currentTime.Add(duration)
		downStops[i+1].Arrival_time = goTime_to_string(currentTime)
		// Add 1-min halt if stop
		if downStops[i+1].IsStop {
			currentTime = currentTime.Add(1 * time.Minute)
		}

		// Departure = arrival + 1 min if IsStop else same as arrival
		downStops[i+1].Departure_time = goTime_to_string(currentTime)

		/*

		   1. Kamaraj College → Kallikudi
		      Duration = 13 min
		      Arrival = 16:40 + 13 min = 16:53
		      Kallikudi IsStop = false → Departure = 16:53

		   2. Kallikudi → Thirumangalam
		      Duration = 14 min
		      Arrival = 16:53 + 14 min = 17:07
		      Thirumangalam IsStop = true → Departure = 17:08

		   3. Thirumangalam → Madurai
		      Duration = 18 min
		      Arrival = 17:08 + 18 min = 17:26
		      Madurai IsStop = true → Departure = 17:27
		*/

		//convert the stop names(normal names )to camel-case to store in DB
		downStops[i].LocationName = Convert_to_CamelCase(downStops[i].LocationName)
	}
	down_route.Stops = downStops
	down_route.ArrivalTime = downStops[len(downStops)-1].Arrival_time

	/*
	   Final Down route table:
	   | Stop            | Arrival | Departure |
	   | --------------- | ------- | --------- |
	   | Kamaraj College | 16:40   | 16:40     |
	   | Kallikudi       | 16:53   | 16:53     |
	   | Thirumangalam   | 17:07   | 17:08     |
	   | Madurai         | 17:26   | 17:27     |
	*/

}
