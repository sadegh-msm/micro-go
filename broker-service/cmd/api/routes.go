package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	c := chi.NewRouter()

	// specify who is allowed to connect
	c.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	c.Get("/gigachad", app.isAlive)
	c.Post("/", app.Broker)
	c.Post("/handle", app.HandleAll)
	c.Post("/log-grpc", app.logEventWithGrpc)
	return c
}
