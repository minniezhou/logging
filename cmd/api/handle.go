package main

import (
	"fmt"
	"net/http"

	"logging-service/cmd/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/minniezhou/jsonToolBox"
)

type Handler struct {
	router *chi.Mux
}

func (c *Config) NewHandler() *Handler {
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

	r.Get("/", c.hello)
	r.Post("/log", c.logToDB)

	return &Handler{
		router: r,
	}
}

func (c *Config) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping Logging Service")

	id, err := c.db.InsertOne("Ping", "Logging for Ping")
	var payLoad jsonToolBox.JsonResponse
	if err != nil {
		payLoad.Error = true
		payLoad.Message = "Ping Logging Service, logging failed"
		fmt.Println(err.Error())
	} else {
		payLoad.Error = false
		payLoad.Message = fmt.Sprintf("Ping Logging Service, Logging Id is %s", id)
	}
	jsonToolBox.WriteJson(w, http.StatusAccepted, payLoad)
}

func (c *Config) logToDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging to DB")
	var logData model.DataType
	err := jsonToolBox.ReadJson(w, r, &logData)
	if err != nil {
		fmt.Println("reading json data error")
		jsonToolBox.ErrorJson(w, err.Error())
		return
	}
	fmt.Println(logData)

	fmt.Println("Inserting one")
	id, err := c.db.InsertOne(logData.Name, logData.Message)
	fmt.Println("Inserted one")
	var payLoad jsonToolBox.JsonResponse
	if err != nil {
		payLoad.Error = true
		payLoad.Message = "logging failed"
		fmt.Println(err.Error())
	} else {
		payLoad.Error = false
		payLoad.Message = fmt.Sprintf("Logging Succeed, Logging Id is %s", id)
		fmt.Println("Inserted one succesfully")
	}

	jsonToolBox.WriteJson(w, http.StatusAccepted, payLoad)
}
