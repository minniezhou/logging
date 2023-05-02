package main

import (
	"context"
	"fmt"
	"log"
	"logging-service/api/logging"
	grpclog "logging-service/cmd/grpc-log"
	"logging-service/cmd/model"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

type Config struct {
	client *mongo.Client
	db     *model.Model
}

const (
	webPort  = 4321
	grpcPort = 43210
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

	go c.grpcListener()

	h := c.NewHandler()
	fmt.Println("starting http server for logging...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", webPort), h.router)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("http logger listening at %d", webPort)
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

func (c *Config) grpcListener() {
	// listen to grpc connection
	fmt.Println("starting grpc listening...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpc_s := grpc.NewServer()
	logServer := grpclog.NewLogServer(c.db)
	logging.RegisterLogServer(grpc_s, logServer)
	fmt.Printf("grpc logger listening at %v", lis.Addr())
	if err := grpc_s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
