package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var pool *pgxpool.Pool

func Connect() {
	var (
		config *pgxpool.Config
		err    error
	)
	//get the connection string from the .env file
	username := os.Getenv("POSTGRES_DB_USERNAME")
	pass := os.Getenv("POSTGRES_DB_PASSWORD")
	port := os.Getenv("POSTGRES_DB_PORT")

	conn := fmt.Sprintf(`postgres://%s:%s@postgres.kws.services:%s/yus?sslmode=disable`, username, pass, port)
	// fmt.Println("conn - ", conn)

	//here context is used for timeout -> if the connection try to connect the db for more than 5 sec , it automatticaly close only that connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // to cancel all the timing after this function ends

	//  ParseConfig parse this "postgres://username:password@host:port/dbname" -> the pool config object so that it can set the maxConnections like that
	config, err = pgxpool.ParseConfig(conn)
	if err != nil {
		fmt.Println("error while parsing the conn - ", err)
	}

	config.MaxConns = 10               // seting 10 maximum connections at the time
	config.MaxConnLifetime = time.Hour //maximum time can the connection live

	pool, err = pgxpool.NewWithConfig(ctx, config) // created the new pool connection with the config
	if err != nil {
		fmt.Println("error while creating the new pool connection - ", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println("error while pool connection ping - ", err)
	}
	fmt.Println("database connected successfully")
}

func GivePostgresConnection() *pgxpool.Pool {
	return pool
}
