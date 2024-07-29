package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Define port on which the server will listen on
const webPort = "8080"

// Define a type, which will be used as a receiver for the application.
type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// connect to rabbitMq
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ!")

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting the broker service on port %s\n", webPort)

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the Server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ is not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
	}

	return connection, nil
}
