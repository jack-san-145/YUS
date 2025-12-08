package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	AppPkg "yus/internal/app"
	"yus/internal/storage"
	"yus/internal/storage/postgres"
	"yus/internal/storage/postgres/service"
	"yus/internal/storage/redis"
)

func main() {

	rc := redis.NewRedisClient()

	AppPkg.App = &AppPkg.Application{
		Port:   "8090",
		Router: NewRouter(),
		Store: &storage.Store{
			InMemoryDB: rc,
		},
	}

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = AppPkg.App.Store.InMemoryDB.CreateClient(ctx) // creates a new redis.Client
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	postgres.Connect() //make a connection to postgres

	go service.Automate_route_scheduling() //change the route direction on-runtime

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", AppPkg.App.Router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
