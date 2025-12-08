package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"yus/internal/storage"
	"yus/internal/storage/postgres"
	"yus/internal/storage/postgres/service"
	"yus/internal/storage/redis"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	type Application struct {
		Port   string
		Router *chi.Mux
		Store  *storage.Store
	}

	rc := redis.NewRedisClient()
	app := &Application{

		Port:   "8090",
		Router: NewRouter(),
		Store: &storage.Store{
			// InMemoryDB: rc,
		},
	}

	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	router := NewRouter()

	postgres.Connect() //make a connection to postgres

	// redis.CreateClient(context.Background()) //made a redis client

	go service.Automate_route_scheduling() //change the route direction on-runtime

	fmt.Println("Server listening on :8090")
	err = http.ListenAndServe("0.0.0.0:8090", router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
