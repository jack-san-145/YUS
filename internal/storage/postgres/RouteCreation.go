package postgres

import (
	"fmt"
	"yus/internal/models"
	"yus/internal/service"
)

func SaveRoute_to_DB(up_route *models.Route) {
	up_route.Direction = "UP"
	down_route := service.Find_down_route(*up_route)
	fmt.Println("down route - ", down_route)
}
