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
	"os"
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
	WEB_PORT            = "4321"
	GRPC_PORT           = "43210"
	MONGO_USER_NAME     = "minniezhou"
	MONGO_USER_PASSWORD = "aJvpaLkL51rUEUlH"
	MONGO_HOST          = "@cluster0.uxdvqoz.mongodb.net/?retryWrites=true&w=majority"
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
	webPort := getEnv("WEB_PORT", WEB_PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%s", webPort), h.router)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("http logger listening at %s", webPort)
}

func connectToMongo() (*mongo.Client, error) {
	userName := getEnv("MONGO_USER_NAME", MONGO_USER_NAME)
	password := getEnv("MONGO_USER_PASSWORD", MONGO_USER_PASSWORD)
	link := getEnv("MONGO_HOST", MONGO_HOST)
	mongo_url := fmt.Sprintf("mongodb+srv://%s:%s%s", userName, password, link)
	log.Println(mongo_url)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_url))
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
	grpcPort := getEnv("GRPC_PORT", GRPC_PORT)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
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

func getEnv(key, default_value string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return default_value
	}
	return value
}
