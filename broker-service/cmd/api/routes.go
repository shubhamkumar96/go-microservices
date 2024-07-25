package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Set up Middleware
	// Specify who is allowed to connect
	mux.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300,
			},
		),
	)

	// To check at this address, if the server is up or not.
	mux.Use(middleware.Heartbeat("/ping"))

	// Add Routes
	mux.Post("/", http.HandlerFunc(app.Broker))

	// Add Route for single point of entry for all the microservice calls
	mux.Post("/handle", app.HandleSubmission)

	return mux
}
