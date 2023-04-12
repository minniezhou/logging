package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbclient *mongo.Client

type Model struct{}

type DataType struct {
	ID      primitive.ObjectID `json:"ID,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`
	Message string             `json:"message" bson:"message,omitempty"`
}

const (
	databaseName   = "Logger"
	collectionName = "log"
)

func NewDBClient(client *mongo.Client) *Model {
	dbclient = client
	return &Model{}
}

func (*Model) InsertOne(name, message string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := DataType{
		Name:    name,
		Message: message,
	}

	collection := dbclient.Database(databaseName).Collection(collectionName)

	result, err := collection.InsertOne(ctx, payload)

	if err != nil {
		return "", err
	}

	fmt.Println(result.InsertedID)
	return result.InsertedID, nil
}
