package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var pool *pgxpool.Pool

type PgStore struct {
	Pool *pgxpool.Pool
}

func NewPgStore() *PgStore {
	return &PgStore{}
}

func (pg *PgStore) Connect(ctx context.Context) error {

	var (
		config *pgxpool.Config
		err    error
	)

	//get the connection string from the .env file
	username := os.Getenv("POSTGRES_DB_USERNAME")
	pass := os.Getenv("POSTGRES_DB_PASSWORD")
	port := os.Getenv("POSTGRES_DB_PORT")

	conn := fmt.Sprintf(`postgres://%s:%s@postgres.kws.services:%s/yus?sslmode=disable`, username, pass, port)

	//  ParseConfig parse this "postgres://username:password@host:port/dbname" -> the pool config object so that it can set the maxConnections like that
	config, err = pgxpool.ParseConfig(conn)
	if err != nil {
		fmt.Println("error while parsing the conn - ", err)
		return err
	}

	config.MaxConns = 10               // seting 10 maximum connections at the time
	config.MaxConnLifetime = time.Hour //maximum time can the connection live

	pg.Pool, err = pgxpool.NewWithConfig(ctx, config) // created the new pool connection with the config
	if err != nil {
		fmt.Println("error while creating the new pool connection - ", err)
		return err
	}

	err = pg.Pool.Ping(ctx)
	if err != nil {
		fmt.Println("error while pool connection ping - ", err)
		return err
	}
	fmt.Println("database connected successfully")
	return nil
}

// func GivePostgresConnection() *pgxpool.Pool {
// 	return pool
// }
