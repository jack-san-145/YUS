package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	AppPkg "yus/internal/AppPkg"
	"yus/internal/handlers"
	"yus/internal/services"
	"yus/internal/storage"
	"yus/internal/storage/postgres"
	"yus/internal/storage/redis"
)

func main() {
	redisStore := redis.NewRedisStore() //Get redisStore
	pgStore := postgres.NewPgStore()

	app := &AppPkg.Application{
		Port: ":8090",
		Store: &storage.Store{
			InMemoryDB: redisStore,
			DB:         pgStore,
		},
	} //Initialize new Application

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Store.InMemoryDB.CreateClient(ctx) // creates a new redis.Client
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	err = app.Store.DB.Connect(ctx) //create a DB connection
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	YUSHandler := handlers.NewHandler(app.Store)

	app.Router = NewRouter(app, YUSHandler)

	app.Store.DB.Connect(ctx)

	go services.AutomateRouteScheduling(app) //to automate route direction change

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe(app.Port, app.Router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
