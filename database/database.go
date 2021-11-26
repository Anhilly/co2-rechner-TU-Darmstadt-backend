package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var client *mongo.Client

func ConnectDatabase() {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	fmt.Println("Connect")

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + username + ":" + password + "@" + serverAdress + ":" + port + "/" + dbName))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect successful")
}

func DisconnectDatabase() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	fmt.Println("Disconnect")
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}