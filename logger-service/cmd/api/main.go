package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
}

func main() {
	// Connect to MongoDB
	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Panicln(err)
	}

	client = mongoClient

	// Create a context in order to disconnect from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close the connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongoDB() (*mongo.Client, error) {
	// Create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: MongoUsername,
		Password: MongoPw,
	})

	// Make connection
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panicln("Error: ", err)
		return nil, err
	}

	return c, err
}
