package models

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Speed     string `json:"speed"`
}

type NominatimResponse struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		Village string `json:"village"`
		Town    string `json:"town"`
		Amenity string `json:"amenity"`
		Road    string `json:"road"`
		City    string `json:"city"`
	} `json:"address"`
}

type Route struct {
	Id            int          `json:"route_id"`
	Name          string       `json:"route_name"`
	Stops         []RouteStops `json:"stops"`
	DepartureTime string       `json:"departure_time"`
	ArrivalTime   string       `json:"arrival_time"`
}

type RouteStops struct {
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	LocationName string `json:"location_name"`
}
