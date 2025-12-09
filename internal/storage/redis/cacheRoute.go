package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"yus/internal/models"
	"yus/internal/storage/postgres"
)

func (r *RedisStore) CacheBusRoute(ctx context.Context) error {
	current_bus_route := postgres.Current_bus_routes()

	current_bus_route_byte, err := json.Marshal(current_bus_route)
	if err != nil {
		fmt.Println("error while marshal the route - ", err)
		return err
	}

	err = r.RedisClient.Set(ctx, "CurrentBusRoute", current_bus_route_byte, time.Hour*24).Err()
	if err != nil {
		fmt.Println("error while set the current_bus_route in redis - ", err)
		return err
	}
	return nil
}

func (r *RedisStore) GetCachedRoute(ctx context.Context) ([]models.CurrentRoute, error) {

	var current_bus_route []models.CurrentRoute

	route_string, err := r.RedisClient.Get(ctx, "CurrentBusRoute").Result()
	if err != nil {
		fmt.Println("error while get the cached bus route - ", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(route_string), &current_bus_route)
	if err != nil {
		fmt.Println("error while unmarshal cached route - ", err)
		return nil, err
	}

	fmt.Println("cached routes from redis ")
	return current_bus_route, nil
}
