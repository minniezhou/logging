package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	router *chi.Mux
}

func NewHandler() *Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Ahuthorization", "Accept", "X-CSRF-Token"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("ping"))

	r.Get("/", hello)
	r.Post("/log", logToDB)

	return &Handler{
		router: r,
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Ping Logging Service")
	payload := jsonRequest{
		Error:   false,
		Message: "Ping Logging Service",
	}
	writeJson(w, http.StatusAccepted, payload)
}

func logToDB(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Logging request to Logging Service")
	payload := jsonRequest{
		Error:   false,
		Message: "Logging Request to Logging Service",
	}
	writeJson(w, http.StatusAccepted, payload)
}
