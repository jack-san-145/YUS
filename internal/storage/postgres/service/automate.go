package service

import (
	"fmt"
	"time"

	"yus/internal/storage/postgres"

	"github.com/robfig/cron/v3"
)

func Automate_route_scheduling() {

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local))

	//for 12 AM
	c.AddFunc("0 0 0 * * *", func() {
		postgres.Change_route_direction("UP")
	})

	//for 12 PM
	c.AddFunc("0 0 12 * * *", func() {
		postgres.Change_route_direction("DOWN")
	})
	c.Start()

	fmt.Println("Route scheduling started..")

	select {} //it blocks the go routiune and run forever
}
