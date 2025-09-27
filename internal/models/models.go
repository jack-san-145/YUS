package models

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
	Dest              string       `json"dest"`
	Stops             []RouteStops `json:"stops"`
	Direction         string       `json:"direction"`
	DepartureTime     string       `json:"up_departure_time"`
	ArrivalTime       string       `json:"up_arrival_time"`
	DownDepartureTime string       `json:"down_departure_time"`
}

type RouteStops struct {
	LocationName   string `json:"location_name"`
	Lat            string `json:"lat"`
	Lon            string `json:"lon"`
	IsStop         bool   `json:"is_stop"`
	Departure_time string `json:"departure_time"`
	Arrival_time   string `json:"arrival_time"`
}
