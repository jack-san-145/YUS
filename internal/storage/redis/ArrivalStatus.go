package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (r *RedisStore) StoreArrivalStatus(ctx context.Context, driverID int, arrivalStatus map[int]string) error {

	// Convert the map to JSON
	arrival_status_json, err := json.Marshal(arrivalStatus)
	if err != nil {
		fmt.Println("error while marshaling arrival status:", err)
		return err
	}

	// Store JSON string in Redis
	err = r.RedisClient.Set(ctx, "ArrivalStatus:"+strconv.Itoa(driverID), arrival_status_json, time.Minute).Err()
	if err != nil {
		fmt.Println("error while storing the arrival status to redis - ", err)
		return err
	}

	return nil
}

func (r *RedisStore) GetArrivalStatus(ctx context.Context, driverID int) (map[int]string, error) {

	var arrival_status_map = make(map[int]string)

	arrival_status_string, err := r.RedisClient.Get(ctx, "ArrivalStatus:"+strconv.Itoa(driverID)).Result()
	if err != nil {
		fmt.Println("error while get arrival status from redis - ", err)
		return arrival_status_map, fmt.Errorf("not found")
	}
	err = json.Unmarshal([]byte(arrival_status_string), &arrival_status_map)
	if err != nil {
		fmt.Println("error while Unmarshaling arrival status:", err)
		return arrival_status_map, fmt.Errorf("not found")
	}
	return arrival_status_map, nil
}
