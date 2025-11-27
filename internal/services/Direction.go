package services

import (
	"time"
)

func Find_current_direction() string {

	current_time := time.Now()
	if current_time.Hour() >= 12 {
		return "DOWN"
	}
	return "UP"
}
