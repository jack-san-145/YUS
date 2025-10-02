package models

import "time"

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Speed     string `json:"speed"`
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
	// DepartureTime string       `json:"down_departure_time"`

}

type RouteStops struct {
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

type AvilableRoutes struct {
	RouteId int    `json:"route_id"`
	Name    string `json:"route_name"`
	Src     string `json:"src"`
	Dest    string `json:"dest"`
	ISUsed  bool   `json:"is_used"` //is any bus using this route currently
	UsedBy  int    `json:"used_by"` // bus using this row now
}
