package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func Store_ArrivalStatus(driver_id int, arrival_status map[int]string) {
	// Convert the map to JSON
	arrival_status_json, err := json.Marshal(arrival_status)
	if err != nil {
		fmt.Println("error while marshaling arrival status:", err)
		return
	}

	// Store JSON string in Redis
	err = rc.Set(context.Background(), "ArrivalStatus:"+strconv.Itoa(driver_id), arrival_status_json, time.Hour).Err()
	if err != nil {
		fmt.Println("error while storing the arrival status to redis - ", err)
		return
	}
}

func Get_ArrivalStatus(driver_id int) (map[int]string, error) {

	var arrival_status_map = make(map[int]string)

	arrival_status_string, err := rc.Get(context.Background(), "ArrivalStatus:"+strconv.Itoa(driver_id)).Result()
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
