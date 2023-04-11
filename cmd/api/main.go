package main

import (
	"context"
	"fmt"
	"log"
	"logging-service/cmd/model"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	client *mongo.Client
	db     *model.Model
}

const (
	webPort = 4321
)

func main() {
	fmt.Println("This is logging service")

	fmt.Println("Connecting to Mongo DB")

	client, err := connectToMongo()
	if err != nil || client == nil {
		fmt.Println("Failed to connect to Mongo DB")
		log.Panic(err)
	}

	fmt.Println("Connected to Mongo DB")
	// connect to MongoDB

	model := model.NewDBClient(client)

	c := Config{
		client: client,
		db:     model,
	}

	h := c.NewHandler()

	err = http.ListenAndServe(fmt.Sprintf(":%d", webPort), h.router)
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	password := url.QueryEscape("aJvpaLkL51rUEUlH")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://minniezhou:" + password + "@cluster0.uxdvqoz.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	//defer client.Disconnect(ctx)
	//defer cancel()
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("ping mongodb failed")
		return nil, err
	}

	return client, nil
}
