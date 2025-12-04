package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"yus/internal/models"
)

func Cache_Bus_Route(current_bus_route []models.CurrentRoute) {
	current_bus_route_byte, err := json.Marshal(current_bus_route)
	if err != nil {
		fmt.Println("error while marshal the route - ", err)
		return
	}

	err = rc.Set(context.Background(), "CurrentBusRoute", current_bus_route_byte, time.Hour*24).Err()
	if err != nil {
		fmt.Println("error while set the current_bus_route in redis - ", err)
		return
	}
}

func Get_cached_route() []models.CurrentRoute {

	var current_bus_route []models.CurrentRoute

	route_string, err := rc.Get(context.Background(), "CurrentBusRoute").Result()
	if err != nil {
		fmt.Println("error while get the cached bus route - ", err)
		return nil
	}
	err = json.Unmarshal([]byte(route_string), &current_bus_route)
	if err != nil {
		fmt.Println("error while unmarshal cached route - ", err)
		return nil
	}

	fmt.Println("cached routes from redis ")
	return current_bus_route
}
