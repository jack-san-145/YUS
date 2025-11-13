package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	// "sync"
	"yus/internal/models"
	"yus/internal/storage/postgres"

	// "yus/internal/service"

	"github.com/gorilla/websocket"
)

func Driver_Ws_hanler(w http.ResponseWriter, r *http.Request) {

	isValid, driver_id := FindDriver_wssSession(r)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("working")
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrading the websocket - ", err)
		return
	}

	Add_Driver_to_passengerMap(driver_id)

	go listen_for_location(driver_id, conn)

}

func listen_for_location(driver_id int, conn *websocket.Conn) {
	defer func() {
		Remove_Driver_from_passengerMap(driver_id)
		conn.Close()
	}()
	fmt.Println("driver connected successfully")

	Ongoing_route := postgres.Find_route_by_busID(driver_id, "DRIVER").Stops
	fmt.Println("ongping route - ", Ongoing_route)

	//bus_status
	// var Arrival_status map[int]string

	// Ping/pong settings
	const (
		pongWait   = 60 * time.Second
		pingPeriod = 50 * time.Second
		writeWait  = 10 * time.Second
	)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start ping ticker
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ping error:", err)
				return
			}
		}
	}()

	var current_location models.Location
	for {
		_, loc, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error while reading the driver's websocket message - ", err)
			return
		}
		err = json.Unmarshal(loc, &current_location)
		if err != nil {
			fmt.Println("error while unmarshaling the location - ", err)
			continue
		}

		Send_location_to_passenger(driver_id, current_location)
		fmt.Printf("latitude - %v & longitude - %v & speed - %v\n",
			current_location.Latitude, current_location.Longitude, current_location.Speed)
	}
}

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
	const threshold = 0.05 // 50 meters

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
	}
	return nearest.StopSequence, reachedTime, is_reached

}
