package main

import (
	"fmt"
	"log"
	"net/http"
)

// Define port on which the server will listen on
const webPort = "80"

// Define a type, which will be used as a receiver for the application.
type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting the broker service on port %s\n", webPort)

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
