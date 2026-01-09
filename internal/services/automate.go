package services

import (
	"context"
	"log"
	"time"

	"yus/internal/AppPkg"

	"github.com/robfig/cron/v3"
)

func AutomateRouteScheduling(app *AppPkg.Application) {

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local))

	//for 12 AM
	c.AddFunc("0 0 0 * * *", func() {

		app.Store.DB.ChangeRouteDirection(context.Background(), "UP")
		current_route, _ := app.Store.DB.GetCurrentBusRoutes(context.Background())
		app.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route) //cache current routes in redis
	})

	//for 12 PM
	c.AddFunc("0 0 12 * * *", func() {
		app.Store.DB.ChangeRouteDirection(context.Background(), "DOWN")
		current_route, _ := app.Store.DB.GetCurrentBusRoutes(context.Background())
		app.Store.InMemoryDB.CacheBusRoute(context.Background(), current_route) //cache current routes in redis
	})
	c.Start()

	log.Println("Route scheduling started..")

	select {} //it blocks the go routiune and run forever
}
