package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	data2 "log-service/data"
	"net"
	"net/http"
	"net/rpc"
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
	Models data2.Models
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

	app := Config{
		Models: data2.New(client),
	}

	// Register the RPC server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	// Start web server
	log.Println("Starting service on port: ", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
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

	log.Println("Connected to mongo")

	return c, err
}
