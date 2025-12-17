package service

import (
	"context"
	"fmt"
	"time"

	"yus/internal/AppPkg"
	"yus/internal/storage/postgres"

	"github.com/robfig/cron/v3"
)

func AutomateRouteScheduling(app *AppPkg.Application) {

	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local))

	//for 12 AM
	c.AddFunc("0 0 0 * * *", func() {
		postgres.ChangeRouteDirection(context.Background(), "UP")
		app.Store.InMemoryDB.CacheBusRoute(context.Background()) //cache current routes in redis
	})

	//for 12 PM
	c.AddFunc("0 0 12 * * *", func() {
		postgres.ChangeRouteDirection(context.Background(), "DOWN")
		app.Store.InMemoryDB.CacheBusRoute(context.Background()) //cache current routes in redis
	})
	c.Start()

	fmt.Println("Route scheduling started..")

	select {} //it blocks the go routiune and run forever
}
