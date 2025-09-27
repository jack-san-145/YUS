package models

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Speed     string `json:"speed"`
}

type Route struct {
	Id                int          `json:"route_id"`
	Name              string       `json:"route_name"`
	Stops             []RouteStops `json:"stops"`
	UpDepartureTime   string       `json:"up_departure_time"`
	UpArrivalTime     string       `json:"up_arrival_time"`
	DownDepartureTime string       `json:"down_departure_time"`
}

type RouteStops struct {
	IsStop       bool   `json:"is_stop"`
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	LocationName string `json:"location_name"`
}
