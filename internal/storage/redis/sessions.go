package redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func generate_sessionID() string {
	return uuid.New().String()
}

func Create_Driver_Session(driver_id int) {
	session_id := generate_sessionID()
	err := rc.HSet(context.Background(), session_id, "category", "DRIVER", "driver_id", driver_id).Err()
	if err != nil {
		fmt.Println("error while setting the new driver session to the redis - ", err)

	}
}
