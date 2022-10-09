package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// CORS policy
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// AVAILABLE ROUTES
	// Provide a route to check live status of the server
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/", app.Broker)
	// Handle all post requests
	mux.Post("/handle", app.HandleSubmission)
	// log to logger service via gRPC
	mux.Post("/log-grpc", app.LogViaGRPC)

	return mux
}
