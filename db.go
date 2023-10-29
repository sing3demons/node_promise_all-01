package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMonoDB() (*mongo.Collection, error) {
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		uri = "mongodb://mongodb1:27017,mongodb2:27018,mongodb3:27019/?replicaSet=my-replica-set"
	}
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URL is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("microservice_db")

	return db.Collection("exampleDb"), nil
}