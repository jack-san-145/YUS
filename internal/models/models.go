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

type RouteStops struct {
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	LocationName string `json:"location_name"`
}
