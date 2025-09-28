package postgres

import (
	"fmt"
	"yus/internal/models"
	"yus/internal/services"
)

func SaveRoute_to_DB(up_route *models.Route) {
	up_route.Direction = "UP"

	services.Calculate_Uproute_departure(up_route)
	down_route := services.Find_down_route(*up_route)
	fmt.Println()
	fmt.Println("up route - ", up_route)
	fmt.Println()
	fmt.Println("down route - ", down_route)
}
