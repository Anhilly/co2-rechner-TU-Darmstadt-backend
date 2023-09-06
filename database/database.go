package database

import (
	"context"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var client *mongo.Client

var (
	dbHost      = ""
	dbPort      = ""
	dbUsername  = ""
	dbPassword  = ""
	dbName      = ""
	dbContainer = ""
)

// ConnectDatabase stellt eine Verbindung mit der Datenbank her mittels der Konstanten aus db_config.go.
// Die Referenz zur Datenbank wird in der Variable client gespeichert
func ConnectDatabase(mode string) error {
	if mode == "dev" { // set values for database connection
		dbHost = config.DevDBHost
		dbPort = config.DevDBPort
		dbUsername = config.DevDBUsername
		dbPassword = config.DevDBPassword
		dbName = config.DevDBName
		dbContainer = config.DevDBContainer
	} else if mode == "prod" {
		dbHost = config.ProdDBHost
		dbPort = config.ProdDBPort
		dbUsername = config.ProdDBUsername
		dbPassword = config.ProdDBPassword
		dbName = config.ProdDBName
		dbContainer = config.ProdDBContainer
	} else {
		log.Fatalln("Mode not specified")
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":" + dbPort))
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

	dir, err := CreateDump("test")
	if err != nil {
		return err
	}
	fmt.Println(dir)

	log.Println("Database Connected")

	return nil
}

// DisconnectDatabase schließt die Verbindung mit der Datenbank.
func DisconnectDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), structs.TimeoutDuration)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("Database Disconnected")

	return nil
}
