package services

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"yus/internal/models"
)

// Haversine returns distance (in km) between two GPS coordinates
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c // distance in km
}

// FindNearestStop finds the stop closest to driver’s current location
func FindNearestStop(driverLat string, driverLon string, stops []models.RouteStops) (int, string, bool) {
	var nearest models.RouteStops
	minDistance := math.MaxFloat64
	const threshold = 0.006 // 6 meters

	driver_lat_float, _ := strconv.ParseFloat(driverLat, 64)
	driver_lon_float, _ := strconv.ParseFloat(driverLon, 64)

	for _, stop := range stops {
		stop_lat_float, _ := strconv.ParseFloat(stop.Lat, 64)
		stop_lon_float, _ := strconv.ParseFloat(stop.Lon, 64)

		distance := Haversine(driver_lat_float, driver_lon_float, stop_lat_float, stop_lon_float)
		if distance < minDistance {
			minDistance = distance
			nearest = stop
		}
	}

	is_reached := minDistance <= threshold
	var reachedTime string
	if is_reached {
		reachedTime = time.Now().Format("15:04")
		fmt.Println("reached time - ", reachedTime)
	}
	return nearest.StopSequence, reachedTime, is_reached

}
