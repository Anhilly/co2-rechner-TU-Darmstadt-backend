package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var client *mongo.Client

/**
Die Funktion stellt eine Verbindung mit der Datenbank her mittels der Konstanten aus db_config.go
Die Referenz zur Datenbank wird in der Variable client gepseichert
*/
func ConnectDatabase() error {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + username + ":" + password + "@" + serverAdress + ":" + port))
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

	return nil
}

/**
Die Funktion schließt die Verbindung mit der Datenbank.
*/
func DisconnectDatabase() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}
