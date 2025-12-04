package service

import (
	"fmt"
	"time"

	"yus/internal/models"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"

	"github.com/robfig/cron/v3"
)

func Automate_route_scheduling() {

	var current_bus_route []models.CurrentRoute

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local))

	//for 12 AM
	c.AddFunc("0 0 0 * * *", func() {
		postgres.Change_route_direction("UP")

		current_bus_route = postgres.Current_bus_routes() //find current bus routes from database
		redis.Cache_Bus_Route(current_bus_route)          //cache current routes in redis
	})

	//for 12 PM
	c.AddFunc("0 0 12 * * *", func() {
		postgres.Change_route_direction("DOWN")

		current_bus_route = postgres.Current_bus_routes()
		redis.Cache_Bus_Route(current_bus_route)
	})
	c.Start()

	fmt.Println("Route scheduling started..")

	select {} //it blocks the go routiune and run forever
}
