package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shubhamkumar96/go-microservices/auth-service/data"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting the Authentication service on port %s\n", webPort)

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Set up Config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the Server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Verify if connecion is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// As we will be adding postgres to docker-compose.yml file,
// so we need to make sure that it is available before we
// actually return the database connection, because this
// service might start before the DB is up.
func connectToDB() *sql.DB {
	// Get 'dsn' string from Env
	dsn := os.Getenv("DSN")

	// for-loop to run for a given time, and to wait till DB
	// connection is established, in which case it breaks out.
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		// to break out of for-loop after some retries.
		if counts > 10 {
			log.Println(err)
			return nil
		}

		// Wait for 2 seconds, before making next call to 'openDB()'
		log.Println("Backing off for 2 Seconds...")
		time.Sleep(2 * time.Second)
	}
}
