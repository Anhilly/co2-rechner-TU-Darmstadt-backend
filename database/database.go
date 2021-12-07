package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var client *mongo.Client

const timeoutDuration time.Duration = 5 * time.Second

/**
Die Funktion stellt eine Verbindung mit der Datenbank her mittels der Konstanten aus db_config.go.
Die Referenz zur Datenbank wird in der Variable client gespeichert
*/
func ConnectDatabase() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	client, err = mongo.NewClient(
		options.Client().ApplyURI("mongodb://" + username + ":" + password + "@" + serverAdress + ":" + port))
	if err != nil {
		return err
	}

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	log.Println("Database Connected")

	return nil
}

/**
Die Funktion schließt die Verbindung mit der Datenbank.
*/
func DisconnectDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("Database Disconnected")

	return nil
}
