package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	AppPkg "yus/internal/AppPkg"
	"yus/internal/handlers"
	"yus/internal/storage"
	"yus/internal/storage/postgres"
	"yus/internal/storage/postgres/service"
	"yus/internal/storage/redis"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	rc := redis.NewRedisClient() //Get redisStore

	app := AppPkg.NewApplication() //Get new Application struct from AppPkg

	app = &AppPkg.Application{
		Port: "8090",
		Store: &storage.Store{
			InMemoryDB: rc,
		},
	} //Initialize new Application

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = app.Store.InMemoryDB.CreateClient(ctx) // creates a new redis.Client
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	YUSHandler := handlers.NewHandler(app.Store)
	app.Router = NewRouter(YUSHandler)

	postgres.Connect() //make a connection to postgres

	go service.Automate_route_scheduling() //change the route direction on-runtime

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", app.Router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
