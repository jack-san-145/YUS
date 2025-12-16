package models

import "time"

type Location struct {
	Latitude      string         `json:"latitude"`
	Longitude     string         `json:"longitude"`
	Speed         string         `json:"speed"`
	ArrivalStatus map[int]string `json:"arrival_status"`
}

type Route struct {
	Id                int          `json:"route_id"`
	UpRouteName       string       `json:"up_route_name"`
	DownRouteName     string       `json:"down_route_name"`
	Src               string       `json:"src"`
	Dest              string       `json:"dest"`
	Stops             []RouteStops `json:"stops"`
	Direction         string       `json:"direction"`
	ArrivalTime       string       `json:"arrival_time"` //its for up_departure_time
	UpDepartureTime   string       `json:"up_departure_time"`
	DownDepartureTime string       `json:"down_departure_time"` //its for down_departure_time
	Created_At        time.Time    `json:"created_at"`
}

type RouteStops struct {
	StopSequence   int    `json:"stop_sequence"`
	LocationName   string `json:"location_name"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
	IsStop         bool   `json:"is_stop"`
	Departure_time string `json:"departure_time"`
	Arrival_time   string `json:"arrival_time"`
}

type Driver struct {
	DriverId  int    `json:"driver_id"`
	Name      string `json:"name"`
	Mobile_no string `json:"mobile_no"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type DriverAddedStatus struct {
	IsAdded  bool   `json:"is_added"`
	Name     string `json:"name"`
	MobileNo string `json:"mobile_no"`
	Email    string `json:"email"`
}

type AvilableRoute struct {
	Id        int    `json:"route_id"`
	Name      string `json:"route_name"`
	Src       string `json:"src"`
	Dest      string `json:"dest"`
	Direction string `json:"direction"`
	Available bool   `json:"available"` //is any bus using this route currently
	BusId     int    `json:"bus_id"`    // bus using this row now
}

type AvailableDriver struct {
	Id        int    `json:"driver_id"`
	Name      string `json:"driver_name"`
	MobileNo  string `json:"mobile_no"`
	Available bool   `json:"available"`
}

type DriverAllocation struct {
	DriverId int `json:"driver_id"`
	BusId    int `json:"bus_id"`
}

type CurrentRoute struct {
	RouteId   int          `json:"route_id"`
	BusId     int          `json:"bus_id"`
	DriverId  int          `json:"driver_id"`
	Direction string       `json:"direction"`
	RouteName string       `json:"route_name"`
	Src       string       `json:"src"`
	Dest      string       `json:"dest"`
	Stops     []RouteStops `json:"stops"`
	IsStop    bool         `json:"is_Stop"`
	Active    bool         `json:"active"`
}

type BusRoute struct {
	BusID     int          `json:"bus_id"`
	RouteId   int          `json:"route_id"`
	RouteName string       `json:"route_name"`
	Src       string       `json:"src"`
	Dest      string       `json:"dest"`
	Stops     []RouteStops `json:"stops"`
	Active    bool         `json:"active"`
}

type AllotedBus struct {
	BusID     int    `json:"bus_id"`
	DriverId  int    `json:"driver_id"`
	RouteId   int    `json:"route_id"`
	RouteName string `json:"route_name"`
	Direction string `json:"direction"`
	Src       string `json:"src"`
	Dest      string `json:"dest"`
}

type CurrentSchedule struct {
	DriverId int `json:"driver_id"`
	BusId    int `json:"bus_id"`
	RouteId  int `json:"route_id"`
}

type PassengerWsRequest struct {
	DriverId  int    `json:"driver_id"`
	RouteId   int    `json:"route_id"`
	Direction string `json:"direction"`
}

type AllRoute struct {
	Currentroute CurrentRoute
	Uproute      CurrentRoute
	Downroute    CurrentRoute
}
